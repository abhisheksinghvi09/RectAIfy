package analyzers

import (
	"context"
	"encoding/json"
	"fmt"

	"rectaify/internal/llm"
	"rectaify/pkg/types"
)

// ExecutionAnalyzer analyzes execution complexity
type ExecutionAnalyzer struct {
	llmClient *llm.Client
}

// NewExecutionAnalyzer creates a new execution analyzer
func NewExecutionAnalyzer(llmClient *llm.Client) *ExecutionAnalyzer {
	return &ExecutionAnalyzer{
		llmClient: llmClient,
	}
}

// Analyze performs execution complexity analysis
func (ea *ExecutionAnalyzer) Analyze(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.ExecutionAnalysis, error) {
	systemPrompt := `You are a startup execution expert. Analyze the provided startup idea and evidence to assess execution complexity.

CRITICAL REQUIREMENTS:
1. ONLY use information explicitly provided in the Evidence
2. If information is not in Evidence, mark as "Unknown" or leave empty
3. Output ONLY valid JSON matching the required schema
4. Reference Evidence by ID numbers when making claims
5. Use exact categories for capital_requirement: "low", "medium", "high", "very high"
6. Use exact categories for talent_rarity: "common", "available", "scarce", "rare"
7. Count integration_count as number of major third-party integrations needed
8. Complexity should be 0.0-1.0 where 1.0 is maximum complexity

Your analysis should focus on:
- Capital requirements based on evidence of similar companies' funding needs
- Talent requirements and availability in the market
- Technical integrations needed (APIs, platforms, services)
- Overall execution complexity combining all factors

Capital requirement guidelines:
- "low": Under $100K, bootstrap-able
- "medium": $100K-$1M, seed round
- "high": $1M-$10M, Series A needed
- "very high": $10M+, multiple rounds

Talent rarity guidelines:
- "common": General business/tech skills
- "available": Specialized but findable skills
- "scarce": Highly specialized, competitive hiring
- "rare": Extremely specialized, very few experts

Base assessments on Evidence, not assumptions.`

	userPrompt := map[string]interface{}{
		"idea":     idea,
		"evidence": evidence,
	}

	schema := []byte(`{
		"type": "object",
		"properties": {
			"capital_requirement": {
				"type": "string",
				"enum": ["low", "medium", "high", "very high"]
			},
			"talent_rarity": {
				"type": "string",
				"enum": ["common", "available", "scarce", "rare"]
			},
			"integration_count": {
				"type": "integer",
				"minimum": 0
			},
			"complexity": {
				"type": "number",
				"minimum": 0.0,
				"maximum": 1.0
			},
			"evidence_ids": {
				"type": "array",
				"items": {"type": "string"}
			}
		},
		"required": ["capital_requirement", "talent_rarity", "integration_count", "complexity", "evidence_ids"],
		"additionalProperties": false
	}`)

	response, err := ea.llmClient.ConstrainedJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return types.ExecutionAnalysis{}, fmt.Errorf("execution analysis failed: %w", err)
	}

	var result types.ExecutionAnalysis
	if err := json.Unmarshal(response, &result); err != nil {
		return types.ExecutionAnalysis{}, fmt.Errorf("failed to parse execution analysis response: %w", err)
	}

	result = ea.validateEvidenceIDs(result, evidence)
	return result, nil
}

func (ea *ExecutionAnalyzer) validateEvidenceIDs(analysis types.ExecutionAnalysis, evidence []types.Evidence) types.ExecutionAnalysis {
	evidenceSet := make(map[string]bool)
	for _, ev := range evidence {
		evidenceSet[ev.ID] = true
	}

	var validEvidenceIDs []string
	for _, id := range analysis.EvidenceIDs {
		if evidenceSet[id] {
			validEvidenceIDs = append(validEvidenceIDs, id)
		}
	}
	analysis.EvidenceIDs = validEvidenceIDs
	return analysis
}
