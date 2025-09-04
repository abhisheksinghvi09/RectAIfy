package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"realitycheck/pkg/types"
)

// Repository handles database operations
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new repository instance
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// SaveAnalysis stores a complete analysis in the database
func (r *Repository) SaveAnalysis(ctx context.Context, analysis types.Analysis) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Marshal idea and result to JSON
	ideaJSON, err := json.Marshal(analysis.Idea)
	if err != nil {
		return fmt.Errorf("failed to marshal idea: %w", err)
	}

	resultJSON, err := json.Marshal(analysis)
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %w", err)
	}

	// Insert analysis
	_, err = tx.Exec(ctx,
		"INSERT INTO analyses (id, idea, result, created_at) VALUES ($1, $2, $3, $4)",
		analysis.ID, ideaJSON, resultJSON, analysis.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert analysis: %w", err)
	}

	// Insert evidence if not already exists and link to analysis
	for _, ev := range analysis.Evidence {
		// Insert evidence (ignore if exists)
		_, err = tx.Exec(ctx,
			`INSERT INTO evidence (id, url, title, snippet, published_at, retrieved_at, source_type) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 ON CONFLICT (id) DO NOTHING`,
			ev.ID, ev.URL, ev.Title, ev.Snippet, ev.PublishedAt, ev.RetrievedAt, ev.SourceType)
		if err != nil {
			return fmt.Errorf("failed to insert evidence %s: %w", ev.ID, err)
		}

		// Link evidence to analysis
		_, err = tx.Exec(ctx,
			`INSERT INTO analysis_evidence (analysis_id, evidence_id) 
			 VALUES ($1, $2)
			 ON CONFLICT DO NOTHING`,
			analysis.ID, ev.ID)
		if err != nil {
			return fmt.Errorf("failed to link evidence %s to analysis %s: %w", ev.ID, analysis.ID, err)
		}
	}

	return tx.Commit(ctx)
}

// GetAnalysis retrieves an analysis by ID
func (r *Repository) GetAnalysis(ctx context.Context, analysisID string) (types.Analysis, error) {
	var resultJSON []byte
	var createdAt time.Time

	err := r.db.QueryRow(ctx,
		"SELECT result, created_at FROM analyses WHERE id = $1",
		analysisID).Scan(&resultJSON, &createdAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return types.Analysis{}, ErrAnalysisNotFound
		}
		return types.Analysis{}, fmt.Errorf("failed to query analysis: %w", err)
	}

	var analysis types.Analysis
	if err := json.Unmarshal(resultJSON, &analysis); err != nil {
		return types.Analysis{}, fmt.Errorf("failed to unmarshal analysis: %w", err)
	}

	// Ensure the timestamps are set correctly
	analysis.CreatedAt = createdAt

	return analysis, nil
}

// GetAnalysisWithEvidence retrieves an analysis with all linked evidence
func (r *Repository) GetAnalysisWithEvidence(ctx context.Context, analysisID string) (types.Analysis, error) {
	analysis, err := r.GetAnalysis(ctx, analysisID)
	if err != nil {
		return types.Analysis{}, err
	}

	// Get linked evidence
	evidence, err := r.GetAnalysisEvidence(ctx, analysisID)
	if err != nil {
		return types.Analysis{}, fmt.Errorf("failed to get analysis evidence: %w", err)
	}

	analysis.Evidence = evidence
	return analysis, nil
}

// GetAnalysisEvidence retrieves all evidence linked to an analysis
func (r *Repository) GetAnalysisEvidence(ctx context.Context, analysisID string) ([]types.Evidence, error) {
	rows, err := r.db.Query(ctx,
		`SELECT e.id, e.url, e.title, e.snippet, e.published_at, e.retrieved_at, e.source_type
		 FROM evidence e
		 JOIN analysis_evidence ae ON e.id = ae.evidence_id
		 WHERE ae.analysis_id = $1
		 ORDER BY e.retrieved_at DESC`,
		analysisID)
	if err != nil {
		return nil, fmt.Errorf("failed to query evidence: %w", err)
	}
	defer rows.Close()

	var evidence []types.Evidence
	for rows.Next() {
		var ev types.Evidence
		err := rows.Scan(&ev.ID, &ev.URL, &ev.Title, &ev.Snippet, &ev.PublishedAt, &ev.RetrievedAt, &ev.SourceType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan evidence: %w", err)
		}
		evidence = append(evidence, ev)
	}

	return evidence, nil
}

// ListAnalyses retrieves a paginated list of analyses
func (r *Repository) ListAnalyses(ctx context.Context, limit, offset int) ([]types.Analysis, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, idea, result, created_at 
		 FROM analyses 
		 ORDER BY created_at DESC 
		 LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query analyses: %w", err)
	}
	defer rows.Close()

	var analyses []types.Analysis
	for rows.Next() {
		var id string
		var ideaJSON, resultJSON []byte
		var createdAt time.Time

		err := rows.Scan(&id, &ideaJSON, &resultJSON, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analysis: %w", err)
		}

		var analysis types.Analysis
		if err := json.Unmarshal(resultJSON, &analysis); err != nil {
			return nil, fmt.Errorf("failed to unmarshal analysis %s: %w", id, err)
		}

		analysis.CreatedAt = createdAt
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// DeleteAnalysis removes an analysis and its evidence links
func (r *Repository) DeleteAnalysis(ctx context.Context, analysisID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete analysis (cascade will handle analysis_evidence)
	result, err := tx.Exec(ctx, "DELETE FROM analyses WHERE id = $1", analysisID)
	if err != nil {
		return fmt.Errorf("failed to delete analysis: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrAnalysisNotFound
	}

	return tx.Commit(ctx)
}

// SaveEvidence stores evidence in the database
func (r *Repository) SaveEvidence(ctx context.Context, evidence []types.Evidence) error {
	if len(evidence) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, ev := range evidence {
		_, err = tx.Exec(ctx,
			`INSERT INTO evidence (id, url, title, snippet, published_at, retrieved_at, source_type) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 ON CONFLICT (id) DO UPDATE SET 
			 url = EXCLUDED.url,
			 title = EXCLUDED.title,
			 snippet = EXCLUDED.snippet,
			 published_at = EXCLUDED.published_at,
			 retrieved_at = EXCLUDED.retrieved_at,
			 source_type = EXCLUDED.source_type`,
			ev.ID, ev.URL, ev.Title, ev.Snippet, ev.PublishedAt, ev.RetrievedAt, ev.SourceType)
		if err != nil {
			return fmt.Errorf("failed to insert evidence %s: %w", ev.ID, err)
		}
	}

	return tx.Commit(ctx)
}

// GetEvidence retrieves evidence by ID
func (r *Repository) GetEvidence(ctx context.Context, evidenceID string) (types.Evidence, error) {
	var ev types.Evidence
	err := r.db.QueryRow(ctx,
		"SELECT id, url, title, snippet, published_at, retrieved_at, source_type FROM evidence WHERE id = $1",
		evidenceID).Scan(&ev.ID, &ev.URL, &ev.Title, &ev.Snippet, &ev.PublishedAt, &ev.RetrievedAt, &ev.SourceType)

	if err != nil {
		if err == pgx.ErrNoRows {
			return types.Evidence{}, ErrEvidenceNotFound
		}
		return types.Evidence{}, fmt.Errorf("failed to query evidence: %w", err)
	}

	return ev, nil
}

// SearchAnalyses searches analyses by idea content
func (r *Repository) SearchAnalyses(ctx context.Context, query string, limit, offset int) ([]types.Analysis, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, idea, result, created_at 
		 FROM analyses 
		 WHERE idea::text ILIKE $1 OR result::text ILIKE $1
		 ORDER BY created_at DESC 
		 LIMIT $2 OFFSET $3`,
		"%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search analyses: %w", err)
	}
	defer rows.Close()

	var analyses []types.Analysis
	for rows.Next() {
		var id string
		var ideaJSON, resultJSON []byte
		var createdAt time.Time

		err := rows.Scan(&id, &ideaJSON, &resultJSON, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analysis: %w", err)
		}

		var analysis types.Analysis
		if err := json.Unmarshal(resultJSON, &analysis); err != nil {
			return nil, fmt.Errorf("failed to unmarshal analysis %s: %w", id, err)
		}

		analysis.CreatedAt = createdAt
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetAnalysisCount returns the total number of analyses
func (r *Repository) GetAnalysisCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM analyses").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count analyses: %w", err)
	}
	return count, nil
}

// CleanupOldEvidence removes evidence older than the specified duration that's not linked to any analysis
func (r *Repository) CleanupOldEvidence(ctx context.Context, olderThan time.Duration) (int, error) {
	cutoff := time.Now().Add(-olderThan)
	
	result, err := r.db.Exec(ctx,
		`DELETE FROM evidence 
		 WHERE retrieved_at < $1 
		 AND id NOT IN (SELECT DISTINCT evidence_id FROM analysis_evidence)`,
		cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old evidence: %w", err)
	}

	return int(result.RowsAffected()), nil
}
