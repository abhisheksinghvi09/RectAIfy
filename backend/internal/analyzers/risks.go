package analyzers

import (
	"context"
	"encoding/json"
	"fmt"

	"realitycheck/internal/llm"
	"realitycheck/pkg/types"
)

// RisksAnalyzer analyzes business risks
type RisksAnalyzer struct {
	llmClient *llm.Client
}

// NewRisksAnalyzer creates a new risks analyzer
func NewRisksAnalyzer(llmClient *llm.Client) *RisksAnalyzer {
	return &RisksAnalyzer{
		llmClient: llmClient,
	}
}

// Analyze performs risk analysis
func (ra *RisksAnalyzer) Analyze(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.RiskAnalysis, error) {
	systemPrompt := `You are a business risk analyst. Analyze the provided startup idea and evidence to identify and assess business risks.

CRITICAL REQUIREMENTS:
1. ONLY use information explicitly provided in the Evidence
2. If information is not in Evidence, mark as "Unknown" or leave empty
3. Output ONLY valid JSON matching the required schema
4. Reference Evidence by ID numbers when making claims
5. Severity and likelihood must be integers 1-5 where 5 is highest/most likely
6. Category should describe the type of risk (e.g., "Market", "Technology", "Financial", "Regulatory")

Your analysis should focus on:
- Market risks: competition, demand changes, market shifts
- Technology risks: technical feasibility, platform dependencies, security
- Financial risks: funding availability, unit economics, cash flow
- Regulatory risks: compliance changes, legal challenges
- Operational risks: talent acquisition, supplier dependencies, execution
- Competitive risks: new entrants, incumbent responses

Risk severity scale (1-5):
1 = Minor impact, easily recoverable
2 = Moderate impact, manageable
3 = Significant impact, requiring major response
4 = Severe impact, threatening business viability
5 = Critical impact, business-ending potential

Risk likelihood scale (1-5):
1 = Very unlikely (< 10% chance)
2 = Unlikely (10-30% chance)
3 = Possible (30-50% chance)
4 = Likely (50-80% chance)
5 = Very likely (> 80% chance)

Only identify risks with Evidence backing. Include mitigation strategies when Evidence suggests them.`

	userPrompt := map[string]interface{}{
		"idea":     idea,
		"evidence": evidence,
	}

	schema := []byte(`{
		"type": "object",
		"properties": {
			"risks": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"category": {"type": "string"},
						"description": {"type": "string"},
						"severity": {
							"type": "integer",
							"minimum": 1,
							"maximum": 5
						},
						"likelihood": {
							"type": "integer",
							"minimum": 1,
							"maximum": 5
						},
						"mitigation": {"type": "string"},
						"evidence_ids": {
							"type": "array",
							"items": {"type": "string"}
						}
					},
					"required": ["category", "description", "severity", "likelihood", "evidence_ids"],
					"additionalProperties": false
				}
			},
			"evidence_ids": {
				"type": "array",
				"items": {"type": "string"}
			}
		},
		"required": ["risks", "evidence_ids"],
		"additionalProperties": false
	}`)

	response, err := ra.llmClient.ConstrainedJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return types.RiskAnalysis{}, fmt.Errorf("risks analysis failed: %w", err)
	}

	var result types.RiskAnalysis
	if err := json.Unmarshal(response, &result); err != nil {
		return types.RiskAnalysis{}, fmt.Errorf("failed to parse risks analysis response: %w", err)
	}

	result = ra.validateEvidenceIDs(result, evidence)
	return result, nil
}

func (ra *RisksAnalyzer) validateEvidenceIDs(analysis types.RiskAnalysis, evidence []types.Evidence) types.RiskAnalysis {
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

	// Validate risk evidence IDs
	for i, risk := range analysis.Risks {
		var validRiskIDs []string
		for _, id := range risk.EvidenceIDs {
			if evidenceSet[id] {
				validRiskIDs = append(validRiskIDs, id)
			}
		}
		analysis.Risks[i].EvidenceIDs = validRiskIDs
	}

	return analysis
}
