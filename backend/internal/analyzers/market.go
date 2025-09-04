package analyzers

import (
	"context"
	"encoding/json"
	"fmt"

	"rectaify/internal/llm"
	"rectaify/pkg/types"
)

// MarketAnalyzer analyzes market conditions and competition
type MarketAnalyzer struct {
	llmClient *llm.Client
}

// NewMarketAnalyzer creates a new market analyzer
func NewMarketAnalyzer(llmClient *llm.Client) *MarketAnalyzer {
	return &MarketAnalyzer{
		llmClient: llmClient,
	}
}

// Analyze performs market analysis based on idea and evidence
func (ma *MarketAnalyzer) Analyze(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.MarketAnalysis, error) {
	// Create the analysis prompt
	systemPrompt := `You are a market research analyst. Analyze the provided startup idea and evidence to assess market conditions.

CRITICAL REQUIREMENTS:
1. ONLY use information explicitly provided in the Evidence
2. If information is not in Evidence, mark as "Unknown" or leave empty
3. Output ONLY valid JSON matching the required schema
4. Reference Evidence by ID numbers when making claims
5. Categorize market stage as exactly one of: "early", "growing", "mature", "declining"

Your analysis should focus on:
- Identifying direct and indirect competitors from Evidence
- Assessing market maturity and growth stage
- Understanding competitive positioning opportunities
- Evaluating market size and opportunity when data is available

Be conservative - if Evidence doesn't clearly support a conclusion, acknowledge uncertainty.`

	// Prepare user prompt with idea and evidence
	userPrompt := map[string]interface{}{
		"idea":     idea,
		"evidence": evidence,
	}

	// Define JSON schema for market analysis
	schema := []byte(`{
		"type": "object",
		"properties": {
			"competitors": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"name": {"type": "string"},
						"description": {"type": "string"},
						"funding": {"type": "string"},
						"stage": {"type": "string"},
						"evidence_ids": {
							"type": "array",
							"items": {"type": "string"}
						}
					},
					"required": ["name", "description", "evidence_ids"],
					"additionalProperties": false
				}
			},
			"market_stage": {
				"type": "string",
				"enum": ["early", "growing", "mature", "declining"]
			},
			"positioning": {"type": "string"},
			"evidence_ids": {
				"type": "array",
				"items": {"type": "string"}
			}
		},
		"required": ["competitors", "market_stage", "positioning", "evidence_ids"],
		"additionalProperties": false
	}`)

	// Call LLM for analysis
	response, err := ma.llmClient.ConstrainedJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return types.MarketAnalysis{}, fmt.Errorf("market analysis failed: %w", err)
	}

	// Parse response
	var result types.MarketAnalysis
	if err := json.Unmarshal(response, &result); err != nil {
		return types.MarketAnalysis{}, fmt.Errorf("failed to parse market analysis response: %w", err)
	}

	// Validate that evidence IDs exist
	result = ma.validateEvidenceIDs(result, evidence)

	return result, nil
}

// validateEvidenceIDs ensures all referenced evidence IDs actually exist
func (ma *MarketAnalyzer) validateEvidenceIDs(analysis types.MarketAnalysis, evidence []types.Evidence) types.MarketAnalysis {
	evidenceSet := make(map[string]bool)
	for _, ev := range evidence {
		evidenceSet[ev.ID] = true
	}

	// Validate main evidence IDs
	var validEvidenceIDs []string
	for _, id := range analysis.EvidenceIDs {
		if evidenceSet[id] {
			validEvidenceIDs = append(validEvidenceIDs, id)
		}
	}
	analysis.EvidenceIDs = validEvidenceIDs

	// Validate competitor evidence IDs
	for i, competitor := range analysis.Competitors {
		var validCompetitorIDs []string
		for _, id := range competitor.EvidenceIDs {
			if evidenceSet[id] {
				validCompetitorIDs = append(validCompetitorIDs, id)
			}
		}
		analysis.Competitors[i].EvidenceIDs = validCompetitorIDs
	}

	return analysis
}
