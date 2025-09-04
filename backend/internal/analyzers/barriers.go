package analyzers

import (
	"context"
	"encoding/json"
	"fmt"

	"realitycheck/internal/llm"
	"realitycheck/pkg/types"
)

// BarriersAnalyzer analyzes execution barriers
type BarriersAnalyzer struct {
	llmClient *llm.Client
}

// NewBarriersAnalyzer creates a new barriers analyzer
func NewBarriersAnalyzer(llmClient *llm.Client) *BarriersAnalyzer {
	return &BarriersAnalyzer{
		llmClient: llmClient,
	}
}

// Analyze performs barrier analysis
func (ba *BarriersAnalyzer) Analyze(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.BarrierAnalysis, error) {
	systemPrompt := `You are a business execution expert. Analyze the provided startup idea and evidence to identify execution barriers.

CRITICAL REQUIREMENTS:
1. ONLY use information explicitly provided in the Evidence
2. If information is not in Evidence, mark as "Unknown" or leave empty
3. Output ONLY valid JSON matching the required schema
4. Reference Evidence by ID numbers when making claims
5. Categorize barriers as exactly one of: "regulation", "supply", "distribution", "trust", "tech"
6. Weight must be between 0.0 and 1.0 representing barrier significance

Your analysis should focus on:
- Regulatory barriers: licensing, compliance, legal requirements, government approval
- Supply barriers: access to materials, suppliers, manufacturing constraints
- Distribution barriers: reaching customers, channel access, logistics
- Trust barriers: building credibility, overcoming skepticism, reputation
- Tech barriers: technical complexity, infrastructure requirements, platform dependencies

Rate barrier weight based on Evidence:
- 0.8-1.0: Major barrier with strong evidence of difficulty
- 0.5-0.7: Moderate barrier with some evidence
- 0.2-0.4: Minor barrier with limited evidence
- 0.0-0.1: Negligible barrier

Be evidence-based - only identify barriers you can substantiate with provided Evidence.`

	userPrompt := map[string]interface{}{
		"idea":     idea,
		"evidence": evidence,
	}

	schema := []byte(`{
		"type": "object",
		"properties": {
			"barriers": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"type": {
							"type": "string",
							"enum": ["regulation", "supply", "distribution", "trust", "tech"]
						},
						"description": {"type": "string"},
						"weight": {
							"type": "number",
							"minimum": 0.0,
							"maximum": 1.0
						},
						"evidence_ids": {
							"type": "array",
							"items": {"type": "string"}
						}
					},
					"required": ["type", "description", "weight", "evidence_ids"],
					"additionalProperties": false
				}
			},
			"evidence_ids": {
				"type": "array",
				"items": {"type": "string"}
			}
		},
		"required": ["barriers", "evidence_ids"],
		"additionalProperties": false
	}`)

	response, err := ba.llmClient.ConstrainedJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return types.BarrierAnalysis{}, fmt.Errorf("barriers analysis failed: %w", err)
	}

	var result types.BarrierAnalysis
	if err := json.Unmarshal(response, &result); err != nil {
		return types.BarrierAnalysis{}, fmt.Errorf("failed to parse barriers analysis response: %w", err)
	}

	result = ba.validateEvidenceIDs(result, evidence)
	return result, nil
}

func (ba *BarriersAnalyzer) validateEvidenceIDs(analysis types.BarrierAnalysis, evidence []types.Evidence) types.BarrierAnalysis {
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

	// Validate barrier evidence IDs
	for i, barrier := range analysis.Barriers {
		var validBarrierIDs []string
		for _, id := range barrier.EvidenceIDs {
			if evidenceSet[id] {
				validBarrierIDs = append(validBarrierIDs, id)
			}
		}
		analysis.Barriers[i].EvidenceIDs = validBarrierIDs
	}

	return analysis
}
