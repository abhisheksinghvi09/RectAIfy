package analyzers

import (
	"context"
	"encoding/json"
	"fmt"

	"realitycheck/internal/llm"
	"realitycheck/internal/score"
	"realitycheck/pkg/types"
)

// VerdictAnalyzer synthesizes all analyses into a final verdict
type VerdictAnalyzer struct {
	llmClient  *llm.Client
	calculator *score.Calculator
}

// NewVerdictAnalyzer creates a new verdict analyzer
func NewVerdictAnalyzer(llmClient *llm.Client, calculator *score.Calculator) *VerdictAnalyzer {
	return &VerdictAnalyzer{
		llmClient:  llmClient,
		calculator: calculator,
	}
}

// Analyze synthesizes all analysis results into a final verdict
func (va *VerdictAnalyzer) Analyze(ctx context.Context, analysis types.Analysis) (types.Viability, error) {
	// First, compute scores using the calculator
	viability := va.calculator.ComputeViability(analysis)

	// Then, enhance with LLM-generated insights
	enhancedViability, err := va.enhanceWithLLMInsights(ctx, analysis, viability)
	if err != nil {
		// If LLM enhancement fails, return the calculated viability
		return viability, nil
	}

	return enhancedViability, nil
}

// enhanceWithLLMInsights adds LLM-generated insights to the computed viability
func (va *VerdictAnalyzer) enhanceWithLLMInsights(ctx context.Context, analysis types.Analysis, viability types.Viability) (types.Viability, error) {
	systemPrompt := `You are a senior startup advisor synthesizing a comprehensive analysis. Review all the analysis components and enhance the verdict with strategic insights.

CRITICAL REQUIREMENTS:
1. ONLY use information from the provided analysis components
2. Output ONLY valid JSON matching the required schema
3. Reference Evidence IDs when making claims
4. DO NOT change the numerical scores - only enhance insights and recommendation
5. Focus on strategic synthesis and actionable insights

Your enhancement should:
- Synthesize insights across all analysis dimensions
- Identify the most critical success/failure factors
- Provide strategic recommendations beyond just the scores
- Highlight key tensions or trade-offs
- Suggest specific next steps for validation or de-risking

Keep insights specific and actionable rather than generic startup advice.`

	userPrompt := map[string]interface{}{
		"analysis":   analysis,
		"viability":  viability,
	}

	schema := []byte(`{
		"type": "object",
		"properties": {
			"overall_score": {"type": "number"},
			"market_score": {"type": "number"},
			"problem_score": {"type": "number"},
			"barrier_score": {"type": "number"},
			"execution_score": {"type": "number"},
			"risk_score": {"type": "number"},
			"graveyard_score": {"type": "number"},
			"recommendation": {"type": "string"},
			"key_insights": {
				"type": "array",
				"items": {"type": "string"}
			},
			"evidence_ids": {
				"type": "array",
				"items": {"type": "string"}
			}
		},
		"required": ["overall_score", "market_score", "problem_score", "barrier_score", "execution_score", "risk_score", "graveyard_score", "recommendation", "key_insights", "evidence_ids"],
		"additionalProperties": false
	}`)

	response, err := va.llmClient.ConstrainedJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return viability, fmt.Errorf("verdict enhancement failed: %w", err)
	}

	var enhancedViability types.Viability
	if err := json.Unmarshal(response, &enhancedViability); err != nil {
		return viability, fmt.Errorf("failed to parse enhanced verdict response: %w", err)
	}

	// Validate evidence IDs
	enhancedViability = va.validateEvidenceIDs(enhancedViability, analysis.Evidence)

	return enhancedViability, nil
}

func (va *VerdictAnalyzer) validateEvidenceIDs(viability types.Viability, evidence []types.Evidence) types.Viability {
	evidenceSet := make(map[string]bool)
	for _, ev := range evidence {
		evidenceSet[ev.ID] = true
	}

	var validEvidenceIDs []string
	for _, id := range viability.EvidenceIDs {
		if evidenceSet[id] {
			validEvidenceIDs = append(validEvidenceIDs, id)
		}
	}
	viability.EvidenceIDs = validEvidenceIDs
	return viability
}
