package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"realitycheck/internal/analyzers"
	"realitycheck/internal/evidence"
	"realitycheck/internal/search"
	"realitycheck/internal/store"
	"realitycheck/pkg/types"
)

// Orchestrator coordinates the entire analysis workflow
type Orchestrator struct {
	planner          *search.Planner
	executor         *search.Executor
	normalizer       *evidence.Normalizer
	coordinator      *analyzers.Coordinator
	repository       *store.Repository
	maxEvidence      int
	analysisTimeout  time.Duration
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(
	planner *search.Planner,
	executor *search.Executor,
	normalizer *evidence.Normalizer,
	coordinator *analyzers.Coordinator,
	repository *store.Repository,
	maxEvidence int,
	analysisTimeout time.Duration,
) *Orchestrator {
	return &Orchestrator{
		planner:         planner,
		executor:        executor,
		normalizer:      normalizer,
		coordinator:     coordinator,
		repository:      repository,
		maxEvidence:     maxEvidence,
		analysisTimeout: analysisTimeout,
	}
}

// AnalyzeIdea performs a complete analysis of a startup idea
func (o *Orchestrator) AnalyzeIdea(ctx context.Context, request types.AnalysisRequest) (string, error) {
	// Create context with timeout
	timeout := o.analysisTimeout
	if request.Options != nil && request.Options.Timeout != nil {
		timeout = *request.Options.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Generate analysis ID
	analysisID, err := o.generateAnalysisID()
	if err != nil {
		return "", fmt.Errorf("failed to generate analysis ID: %w", err)
	}

	// Step 1: Plan search queries
	queries, err := o.planner.Plan(ctx, request.Idea)
	if err != nil {
		return "", fmt.Errorf("query planning failed: %w", err)
	}

	// Step 2: Execute searches and gather evidence
	location := request.Options.GetLocation()
	rawEvidence, err := o.executor.Run(ctx, queries, location)
	if err != nil {
		return "", fmt.Errorf("search execution failed: %w", err)
	}

	// Step 3: Normalize and deduplicate evidence
	normalizedEvidence := o.normalizer.Normalize(rawEvidence)

	// Step 4: Limit evidence if needed
	maxEvidence := o.maxEvidence
	if request.Options != nil && request.Options.MaxEvidence > 0 {
		maxEvidence = request.Options.MaxEvidence
	}
	if len(normalizedEvidence) > maxEvidence {
		normalizedEvidence = normalizedEvidence[:maxEvidence]
	}

	// Step 5: Run all analyzers
	analysis, err := o.coordinator.AnalyzeAll(ctx, request.Idea, normalizedEvidence)
	if err != nil {
		return "", fmt.Errorf("analysis failed: %w", err)
	}

	// Step 6: Finalize analysis metadata
	analysis.ID = analysisID
	analysis.CreatedAt = time.Now()

	// Check if context was cancelled (partial analysis)
	select {
	case <-ctx.Done():
		analysis.Partial = true
	default:
	}

	// Step 7: Save to database
	if err := o.repository.SaveAnalysis(ctx, analysis); err != nil {
		return "", fmt.Errorf("failed to save analysis: %w", err)
	}

	return analysisID, nil
}

// GetAnalysis retrieves a stored analysis
func (o *Orchestrator) GetAnalysis(ctx context.Context, analysisID string) (types.Analysis, error) {
	return o.repository.GetAnalysisWithEvidence(ctx, analysisID)
}

// ListAnalyses returns a paginated list of analyses
func (o *Orchestrator) ListAnalyses(ctx context.Context, limit, offset int) ([]types.Analysis, error) {
	return o.repository.ListAnalyses(ctx, limit, offset)
}

// SearchAnalyses searches for analyses matching a query
func (o *Orchestrator) SearchAnalyses(ctx context.Context, query string, limit, offset int) ([]types.Analysis, error) {
	return o.repository.SearchAnalyses(ctx, query, limit, offset)
}

// DeleteAnalysis removes an analysis
func (o *Orchestrator) DeleteAnalysis(ctx context.Context, analysisID string) error {
	return o.repository.DeleteAnalysis(ctx, analysisID)
}

// GetAnalysisCount returns the total number of analyses
func (o *Orchestrator) GetAnalysisCount(ctx context.Context) (int, error) {
	return o.repository.GetAnalysisCount(ctx)
}

// generateAnalysisID creates a unique analysis identifier
func (o *Orchestrator) generateAnalysisID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HealthCheck performs a basic health check of all components
func (o *Orchestrator) HealthCheck(ctx context.Context) error {
	// Check database connectivity
	count, err := o.repository.GetAnalysisCount(ctx)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Basic validation that we can access the database
	_ = count

	return nil
}

// GetStats returns basic statistics about the system
func (o *Orchestrator) GetStats(ctx context.Context) (map[string]interface{}, error) {
	totalAnalyses, err := o.repository.GetAnalysisCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis count: %w", err)
	}

	stats := map[string]interface{}{
		"total_analyses": totalAnalyses,
		"max_evidence":   o.maxEvidence,
		"timeout":        o.analysisTimeout.String(),
	}

	return stats, nil
}

// CleanupOldData removes old evidence that's not linked to analyses
func (o *Orchestrator) CleanupOldData(ctx context.Context, olderThan time.Duration) (int, error) {
	return o.repository.CleanupOldEvidence(ctx, olderThan)
}
