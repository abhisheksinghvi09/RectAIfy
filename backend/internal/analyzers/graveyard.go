package analyzers

import (
	"context"
	"encoding/json"
	"fmt"

	"rectaify/internal/llm"
	"rectaify/pkg/types"
)

// GraveyardAnalyzer analyzes failed similar companies
type GraveyardAnalyzer struct {
	llmClient *llm.Client
}

// NewGraveyardAnalyzer creates a new graveyard analyzer
func NewGraveyardAnalyzer(llmClient *llm.Client) *GraveyardAnalyzer {
	return &GraveyardAnalyzer{
		llmClient: llmClient,
	}
}

// Analyze performs graveyard analysis
func (ga *GraveyardAnalyzer) Analyze(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.GraveyardAnalysis, error) {
	systemPrompt := `You are a startup postmortem analyst. Analyze the provided startup idea and evidence to identify failed similar companies and extract lessons.

CRITICAL REQUIREMENTS:
1. ONLY use information explicitly provided in the Evidence
2. If information is not in Evidence, mark as "Unknown" or leave empty
3. Output ONLY valid JSON matching the required schema
4. Reference Evidence by ID numbers when making claims
5. Only include companies with clear evidence of failure/shutdown
6. Focus on extracting actionable lessons from failures

Your analysis should focus on:
- Companies that attempted similar solutions and failed
- Clear failure causes backed by evidence (not speculation)
- Specific lessons that can be learned from each failure
- Patterns across multiple failures if present

Types of failure causes to look for:
- Market: No demand, wrong timing, market too small
- Product: Poor execution, technical issues, bad UX
- Business model: Unsustainable economics, pricing issues
- Competition: Outcompeted, market consolidated
- Funding: Couldn't raise capital, burned through money
- Team: Founder issues, key departures, execution problems
- External: Regulatory changes, economic conditions

Extract specific, actionable lessons rather than generic advice. Only include cases with solid evidence backing.`

	userPrompt := map[string]interface{}{
		"idea":     idea,
		"evidence": evidence,
	}

	schema := []byte(`{
		"type": "object",
		"properties": {
			"cases": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"company_name": {"type": "string"},
						"description": {"type": "string"},
						"failure_cause": {"type": "string"},
						"lessons": {"type": "string"},
						"evidence_ids": {
							"type": "array",
							"items": {"type": "string"}
						}
					},
					"required": ["company_name", "description", "failure_cause", "lessons", "evidence_ids"],
					"additionalProperties": false
				}
			},
			"evidence_ids": {
				"type": "array",
				"items": {"type": "string"}
			}
		},
		"required": ["cases", "evidence_ids"],
		"additionalProperties": false
	}`)

	response, err := ga.llmClient.ConstrainedJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return types.GraveyardAnalysis{}, fmt.Errorf("graveyard analysis failed: %w", err)
	}

	var result types.GraveyardAnalysis
	if err := json.Unmarshal(response, &result); err != nil {
		return types.GraveyardAnalysis{}, fmt.Errorf("failed to parse graveyard analysis response: %w", err)
	}

	result = ga.validateEvidenceIDs(result, evidence)
	return result, nil
}

func (ga *GraveyardAnalyzer) validateEvidenceIDs(analysis types.GraveyardAnalysis, evidence []types.Evidence) types.GraveyardAnalysis {
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

	// Validate case evidence IDs
	for i, graveyardCase := range analysis.Cases {
		var validCaseIDs []string
		for _, id := range graveyardCase.EvidenceIDs {
			if evidenceSet[id] {
				validCaseIDs = append(validCaseIDs, id)
			}
		}
		analysis.Cases[i].EvidenceIDs = validCaseIDs
	}

	return analysis
}
