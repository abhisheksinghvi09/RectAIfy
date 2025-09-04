package search

import (
	"context"
	"fmt"
	"sync"
	"time"

	"rectaify/internal/cache"
	"rectaify/internal/llm"
	"rectaify/pkg/types"
)

// Executor handles search query execution with caching
type Executor struct {
	llmClient *llm.Client
	cache     *cache.EvidenceCache
	timeout   time.Duration
}

// NewExecutor creates a new search executor
func NewExecutor(llmClient *llm.Client, evidenceCache *cache.EvidenceCache, timeout time.Duration) *Executor {
	return &Executor{
		llmClient: llmClient,
		cache:     evidenceCache,
		timeout:   timeout,
	}
}

// Run executes a batch of search queries with caching and deduplication
func (e *Executor) Run(ctx context.Context, queries []types.SearchQuery, location *types.ApproxLocation) ([]types.Evidence, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Group queries by priority and process in batches
	batches := e.groupQueriesByPriority(queries)
	
	var allEvidence []types.Evidence
	var mu sync.Mutex
	
	// Process each priority batch
	for priority := 1; priority <= 3; priority++ {
		if priorityQueries, exists := batches[priority]; exists {
			evidence, err := e.processBatch(ctx, priorityQueries, location)
			if err != nil {
				// Log error but continue with other batches
				continue
			}
			
			mu.Lock()
			allEvidence = append(allEvidence, evidence...)
			mu.Unlock()
		}
	}
	
	// Deduplicate evidence
	deduped := e.deduplicateEvidence(allEvidence)
	
	return deduped, nil
}

// processBatch processes a batch of queries with the same priority
func (e *Executor) processBatch(ctx context.Context, queries []types.SearchQuery, location *types.ApproxLocation) ([]types.Evidence, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allEvidence []types.Evidence
	
	// Limit concurrent searches
	sem := make(chan struct{}, 3) // Max 3 concurrent searches
	
	for _, query := range queries {
		wg.Add(1)
		
		go func(q types.SearchQuery) {
			defer wg.Done()
			
			// Acquire semaphore
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return
			}
			
			evidence, err := e.executeQuery(ctx, q, location)
			if err != nil {
				// Log error but continue
				return
			}
			
			mu.Lock()
			allEvidence = append(allEvidence, evidence...)
			mu.Unlock()
		}(query)
	}
	
	wg.Wait()
	return allEvidence, nil
}

// executeQuery executes a single search query with caching
func (e *Executor) executeQuery(ctx context.Context, query types.SearchQuery, location *types.ApproxLocation) ([]types.Evidence, error) {
	// Create cache key that includes location context
	cacheKey := e.createCacheKey(query.Query, location)
	
	// Check cache first
	if cached, found, err := e.cache.GetEvidence(ctx, cacheKey); err == nil && found {
		return cached, nil
	}
	
	// Execute search via LLM client
	evidence, err := e.llmClient.Search(ctx, []string{query.Query}, location)
	if err != nil {
		return nil, fmt.Errorf("search failed for query '%s': %w", query.Query, err)
	}
	
	// Store in cache
	if err := e.cache.SetEvidence(ctx, cacheKey, evidence); err != nil {
		// Log cache error but don't fail the request
	}
	
	return evidence, nil
}

// groupQueriesByPriority groups queries by their priority level
func (e *Executor) groupQueriesByPriority(queries []types.SearchQuery) map[int][]types.SearchQuery {
	batches := make(map[int][]types.SearchQuery)
	
	for _, query := range queries {
		priority := query.Priority
		if priority < 1 || priority > 3 {
			priority = 3 // Default to lowest priority
		}
		batches[priority] = append(batches[priority], query)
	}
	
	return batches
}

// createCacheKey creates a cache key that includes location context
func (e *Executor) createCacheKey(query string, location *types.ApproxLocation) string {
	key := query
	
	if location != nil {
		if location.Country != "" {
			key += "|country:" + location.Country
		}
		if location.Region != "" {
			key += "|region:" + location.Region
		}
	}
	
	return key
}

// deduplicateEvidence removes duplicate evidence entries
func (e *Executor) deduplicateEvidence(evidence []types.Evidence) []types.Evidence {
	seen := make(map[string]bool)
	var unique []types.Evidence
	
	for _, ev := range evidence {
		// Use URL + title as deduplication key
		key := ev.URL + "|" + ev.Title
		
		if !seen[key] {
			seen[key] = true
			unique = append(unique, ev)
		}
	}
	
	return unique
}
