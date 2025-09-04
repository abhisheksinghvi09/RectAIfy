package analyzers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"realitycheck/internal/llm"
	"realitycheck/internal/score"
	"realitycheck/pkg/types"
)

// Coordinator manages all analyzers and runs them in parallel
type Coordinator struct {
	marketAnalyzer     *MarketAnalyzer
	problemAnalyzer    *ProblemAnalyzer
	barriersAnalyzer   *BarriersAnalyzer
	executionAnalyzer  *ExecutionAnalyzer
	risksAnalyzer      *RisksAnalyzer
	graveyardAnalyzer  *GraveyardAnalyzer
	verdictAnalyzer    *VerdictAnalyzer
}

// NewCoordinator creates a new analyzer coordinator
func NewCoordinator(llmClient *llm.Client, calculator *score.Calculator) *Coordinator {
	return &Coordinator{
		marketAnalyzer:     NewMarketAnalyzer(llmClient),
		problemAnalyzer:    NewProblemAnalyzer(llmClient),
		barriersAnalyzer:   NewBarriersAnalyzer(llmClient),
		executionAnalyzer:  NewExecutionAnalyzer(llmClient),
		risksAnalyzer:      NewRisksAnalyzer(llmClient),
		graveyardAnalyzer:  NewGraveyardAnalyzer(llmClient),
		verdictAnalyzer:    NewVerdictAnalyzer(llmClient, calculator),
	}
}

// AnalyzeAll runs all analyzers in parallel and returns complete analysis
func (c *Coordinator) AnalyzeAll(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.Analysis, error) {
	// Run all analyzers in parallel except verdict (which depends on others)
	var market types.MarketAnalysis
	var problem types.ProblemAnalysis
	var barriers types.BarrierAnalysis
	var execution types.ExecutionAnalysis
	var risks types.RiskAnalysis
	var graveyard types.GraveyardAnalysis

	var mu sync.Mutex
	var analysisErrors []error

	g, ctx := errgroup.WithContext(ctx)

	// Market analysis
	g.Go(func() error {
		result, err := c.marketAnalyzer.Analyze(ctx, idea, evidence)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("market analysis failed: %w", err))
			mu.Unlock()
			return nil // Don't fail the entire group
		}
		mu.Lock()
		market = result
		mu.Unlock()
		return nil
	})

	// Problem analysis
	g.Go(func() error {
		result, err := c.problemAnalyzer.Analyze(ctx, idea, evidence)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("problem analysis failed: %w", err))
			mu.Unlock()
			return nil
		}
		mu.Lock()
		problem = result
		mu.Unlock()
		return nil
	})

	// Barriers analysis
	g.Go(func() error {
		result, err := c.barriersAnalyzer.Analyze(ctx, idea, evidence)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("barriers analysis failed: %w", err))
			mu.Unlock()
			return nil
		}
		mu.Lock()
		barriers = result
		mu.Unlock()
		return nil
	})

	// Execution analysis
	g.Go(func() error {
		result, err := c.executionAnalyzer.Analyze(ctx, idea, evidence)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("execution analysis failed: %w", err))
			mu.Unlock()
			return nil
		}
		mu.Lock()
		execution = result
		mu.Unlock()
		return nil
	})

	// Risks analysis
	g.Go(func() error {
		result, err := c.risksAnalyzer.Analyze(ctx, idea, evidence)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("risks analysis failed: %w", err))
			mu.Unlock()
			return nil
		}
		mu.Lock()
		risks = result
		mu.Unlock()
		return nil
	})

	// Graveyard analysis
	g.Go(func() error {
		result, err := c.graveyardAnalyzer.Analyze(ctx, idea, evidence)
		if err != nil {
			mu.Lock()
			analysisErrors = append(analysisErrors, fmt.Errorf("graveyard analysis failed: %w", err))
			mu.Unlock()
			return nil
		}
		mu.Lock()
		graveyard = result
		mu.Unlock()
		return nil
	})

	// Wait for all analyzers to complete
	if err := g.Wait(); err != nil {
		return types.Analysis{}, err
	}

	// Create preliminary analysis for verdict
	preliminaryAnalysis := types.Analysis{
		Idea:      idea,
		Market:    market,
		Problem:   problem,
		Barriers:  barriers,
		Execution: execution,
		Risks:     risks,
		Graveyard: graveyard,
		Evidence:  evidence,
	}

	// Run verdict analysis
	verdict, err := c.verdictAnalyzer.Analyze(ctx, preliminaryAnalysis)
	if err != nil {
		analysisErrors = append(analysisErrors, fmt.Errorf("verdict analysis failed: %w", err))
		// Use empty verdict if it fails
		verdict = types.Viability{}
	}

	// Final analysis
	finalAnalysis := types.Analysis{
		Idea:      idea,
		Market:    market,
		Problem:   problem,
		Barriers:  barriers,
		Execution: execution,
		Risks:     risks,
		Graveyard: graveyard,
		Verdict:   verdict,
		Evidence:  evidence,
		Partial:   len(analysisErrors) > 0,
	}

	// Include error information in meta if there were issues
	if len(analysisErrors) > 0 {
		errorMeta := map[string]interface{}{
			"errors": analysisErrors,
		}
		if metaBytes, err := json.Marshal(errorMeta); err == nil {
			finalAnalysis.Meta = metaBytes
		}
	}

	return finalAnalysis, nil
}

// AnalyzeMarket runs only market analysis (for testing/debugging)
func (c *Coordinator) AnalyzeMarket(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.MarketAnalysis, error) {
	return c.marketAnalyzer.Analyze(ctx, idea, evidence)
}

// AnalyzeProblem runs only problem analysis (for testing/debugging)
func (c *Coordinator) AnalyzeProblem(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.ProblemAnalysis, error) {
	return c.problemAnalyzer.Analyze(ctx, idea, evidence)
}

// AnalyzeBarriers runs only barriers analysis (for testing/debugging)
func (c *Coordinator) AnalyzeBarriers(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.BarrierAnalysis, error) {
	return c.barriersAnalyzer.Analyze(ctx, idea, evidence)
}

// AnalyzeExecution runs only execution analysis (for testing/debugging)
func (c *Coordinator) AnalyzeExecution(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.ExecutionAnalysis, error) {
	return c.executionAnalyzer.Analyze(ctx, idea, evidence)
}

// AnalyzeRisks runs only risks analysis (for testing/debugging)
func (c *Coordinator) AnalyzeRisks(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.RiskAnalysis, error) {
	return c.risksAnalyzer.Analyze(ctx, idea, evidence)
}

// AnalyzeGraveyard runs only graveyard analysis (for testing/debugging)
func (c *Coordinator) AnalyzeGraveyard(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.GraveyardAnalysis, error) {
	return c.graveyardAnalyzer.Analyze(ctx, idea, evidence)
}
