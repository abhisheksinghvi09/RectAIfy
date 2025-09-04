package analyzers

import (
	"context"
	"encoding/json"
	"fmt"

	"rectaify/internal/llm"
	"rectaify/pkg/types"
)

// ProblemAnalyzer analyzes problem validation and pain points
type ProblemAnalyzer struct {
	llmClient *llm.Client
}

// NewProblemAnalyzer creates a new problem analyzer
func NewProblemAnalyzer(llmClient *llm.Client) *ProblemAnalyzer {
	return &ProblemAnalyzer{
		llmClient: llmClient,
	}
}

// Analyze performs problem validation analysis
func (pa *ProblemAnalyzer) Analyze(ctx context.Context, idea types.IdeaInput, evidence []types.Evidence) (types.ProblemAnalysis, error) {
	systemPrompt := `You are a problem validation expert. Analyze the provided startup idea and evidence to assess problem validity.

CRITICAL REQUIREMENTS:
1. ONLY use information explicitly provided in the Evidence
2. If information is not in Evidence, mark as "Unknown" or leave empty
3. Output ONLY valid JSON matching the required schema
4. Reference Evidence by ID numbers when making claims
5. Focus on identifying real user pain points and validation signals

Your analysis should focus on:
- Identifying specific pain points that users actually experience
- Finding evidence of user complaints, frustrations, or current workarounds
- Assessing whether the problem is widespread vs niche
- Evaluating problem urgency and frequency
- Looking for validation signals like user-generated content, forum discussions, surveys

Be skeptical - distinguish between assumed problems and evidence-backed pain points.`

	userPrompt := map[string]interface{}{
		"idea":     idea,
		"evidence": evidence,
	}

	schema := []byte(`{
		"type": "object",
		"properties": {
			"pain_points": {
				"type": "array",
				"items": {"type": "string"},
				"description": "Specific pain points with evidence backing"
			},
			"validation": {
				"type": "string",
				"description": "Summary of problem validation evidence"
			},
			"evidence_ids": {
				"type": "array",
				"items": {"type": "string"}
			}
		},
		"required": ["pain_points", "validation", "evidence_ids"],
		"additionalProperties": false
	}`)

	response, err := pa.llmClient.ConstrainedJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return types.ProblemAnalysis{}, fmt.Errorf("problem analysis failed: %w", err)
	}

	var result types.ProblemAnalysis
	if err := json.Unmarshal(response, &result); err != nil {
		return types.ProblemAnalysis{}, fmt.Errorf("failed to parse problem analysis response: %w", err)
	}

	result = pa.validateEvidenceIDs(result, evidence)
	return result, nil
}

func (pa *ProblemAnalyzer) validateEvidenceIDs(analysis types.ProblemAnalysis, evidence []types.Evidence) types.ProblemAnalysis {
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
