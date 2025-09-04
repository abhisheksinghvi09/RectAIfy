package report

import (
	"fmt"
	"strings"

	"rectaify/pkg/types"
)

// MarkdownBuilder generates markdown reports from analysis results
type MarkdownBuilder struct{}

// NewMarkdownBuilder creates a new markdown builder
func NewMarkdownBuilder() *MarkdownBuilder {
	return &MarkdownBuilder{}
}

// Build generates a markdown report from analysis
func (mb *MarkdownBuilder) Build(analysis types.Analysis) string {
	var report strings.Builder

	// Header
	report.WriteString(fmt.Sprintf("# RectAify: %s\n\n", analysis.Idea.Title))
	report.WriteString(fmt.Sprintf("**One-liner:** %s\n\n", analysis.Idea.OneLiner))
	report.WriteString(fmt.Sprintf("**Analysis Date:** %s\n\n", analysis.CreatedAt.Format("January 2, 2006")))

	if analysis.Partial {
		report.WriteString("⚠️ **Note:** This analysis is partial due to timeout or processing limitations.\n\n")
	}

	// Executive Summary
	report.WriteString("## Executive Summary\n\n")
	report.WriteString(fmt.Sprintf("**Overall Score:** %.1f/100\n\n", analysis.Verdict.OverallScore))
	report.WriteString(fmt.Sprintf("**Recommendation:** %s\n\n", analysis.Verdict.Recommendation))

	// Score Breakdown
	report.WriteString("### Score Breakdown\n\n")
	report.WriteString("| Dimension | Score | Assessment |\n")
	report.WriteString("|-----------|-------|------------|\n")
	report.WriteString(fmt.Sprintf("| Market | %.1f/100 | %s |\n", analysis.Verdict.MarketScore, mb.getScoreAssessment(analysis.Verdict.MarketScore)))
	report.WriteString(fmt.Sprintf("| Problem | %.1f/100 | %s |\n", analysis.Verdict.ProblemScore, mb.getScoreAssessment(analysis.Verdict.ProblemScore)))
	report.WriteString(fmt.Sprintf("| Barriers | %.1f/100 | %s |\n", analysis.Verdict.BarrierScore, mb.getScoreAssessment(analysis.Verdict.BarrierScore)))
	report.WriteString(fmt.Sprintf("| Execution | %.1f/100 | %s |\n", analysis.Verdict.ExecutionScore, mb.getScoreAssessment(analysis.Verdict.ExecutionScore)))
	report.WriteString(fmt.Sprintf("| Risks | %.1f/100 | %s |\n", analysis.Verdict.RiskScore, mb.getScoreAssessment(analysis.Verdict.RiskScore)))
	report.WriteString(fmt.Sprintf("| Graveyard | %.1f/100 | %s |\n", analysis.Verdict.GraveyardScore, mb.getScoreAssessment(analysis.Verdict.GraveyardScore)))
	report.WriteString("\n")

	// Key Insights
	if len(analysis.Verdict.KeyInsights) > 0 {
		report.WriteString("### Key Insights\n\n")
		for _, insight := range analysis.Verdict.KeyInsights {
			report.WriteString(fmt.Sprintf("- %s\n", insight))
		}
		report.WriteString("\n")
	}

	// Detailed Analysis
	report.WriteString("## Detailed Analysis\n\n")

	// Market Analysis
	report.WriteString("### Market Analysis\n\n")
	report.WriteString(fmt.Sprintf("**Market Stage:** %s\n\n", strings.Title(analysis.Market.MarketStage)))
	if analysis.Market.Positioning != "" {
		report.WriteString(fmt.Sprintf("**Positioning:** %s\n\n", analysis.Market.Positioning))
	}

	if len(analysis.Market.Competitors) > 0 {
		report.WriteString("#### Competitors\n\n")
		for i, competitor := range analysis.Market.Competitors {
			report.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, competitor.Name))
			report.WriteString(fmt.Sprintf("   - %s\n", competitor.Description))
			if competitor.Funding != "" {
				report.WriteString(fmt.Sprintf("   - Funding: %s\n", competitor.Funding))
			}
			if competitor.Stage != "" {
				report.WriteString(fmt.Sprintf("   - Stage: %s\n", competitor.Stage))
			}
			if len(competitor.EvidenceIDs) > 0 {
				report.WriteString(fmt.Sprintf("   - Sources: %s\n", mb.formatEvidenceRefs(competitor.EvidenceIDs)))
			}
			report.WriteString("\n")
		}
	}

	// Problem Analysis
	report.WriteString("### Problem Analysis\n\n")
	if len(analysis.Problem.PainPoints) > 0 {
		report.WriteString("#### Pain Points\n\n")
		for i, painPoint := range analysis.Problem.PainPoints {
			report.WriteString(fmt.Sprintf("%d. %s\n", i+1, painPoint))
		}
		report.WriteString("\n")
	}

	if analysis.Problem.Validation != "" {
		report.WriteString("#### Validation\n\n")
		report.WriteString(fmt.Sprintf("%s\n\n", analysis.Problem.Validation))
	}

	// Barriers Analysis
	if len(analysis.Barriers.Barriers) > 0 {
		report.WriteString("### Execution Barriers\n\n")
		for i, barrier := range analysis.Barriers.Barriers {
			weight := barrier.Weight * 100
			report.WriteString(fmt.Sprintf("%d. **%s** (Impact: %.0f%%)\n", i+1, strings.Title(barrier.Type), weight))
			report.WriteString(fmt.Sprintf("   %s\n", barrier.Description))
			if len(barrier.EvidenceIDs) > 0 {
				report.WriteString(fmt.Sprintf("   Sources: %s\n", mb.formatEvidenceRefs(barrier.EvidenceIDs)))
			}
			report.WriteString("\n")
		}
	}

	// Execution Analysis
	report.WriteString("### Execution Analysis\n\n")
	report.WriteString(fmt.Sprintf("**Capital Requirement:** %s\n", strings.Title(analysis.Execution.CapitalRequirement)))
	report.WriteString(fmt.Sprintf("**Talent Rarity:** %s\n", strings.Title(analysis.Execution.TalentRarity)))
	report.WriteString(fmt.Sprintf("**Integration Count:** %d\n", analysis.Execution.IntegrationCount))
	report.WriteString(fmt.Sprintf("**Complexity Score:** %.2f/1.0\n\n", analysis.Execution.Complexity))

	// Risk Analysis
	if len(analysis.Risks.Risks) > 0 {
		report.WriteString("### Risk Analysis\n\n")
		for i, risk := range analysis.Risks.Risks {
			impact := risk.Severity * risk.Likelihood
			report.WriteString(fmt.Sprintf("%d. **%s Risk** (Severity: %d/5, Likelihood: %d/5, Impact: %d/25)\n", 
				i+1, risk.Category, risk.Severity, risk.Likelihood, impact))
			report.WriteString(fmt.Sprintf("   %s\n", risk.Description))
			if risk.Mitigation != "" {
				report.WriteString(fmt.Sprintf("   **Mitigation:** %s\n", risk.Mitigation))
			}
			if len(risk.EvidenceIDs) > 0 {
				report.WriteString(fmt.Sprintf("   Sources: %s\n", mb.formatEvidenceRefs(risk.EvidenceIDs)))
			}
			report.WriteString("\n")
		}
	}

	// Graveyard Analysis
	if len(analysis.Graveyard.Cases) > 0 {
		report.WriteString("### Graveyard Analysis\n\n")
		report.WriteString("#### Failed Similar Companies\n\n")
		for i, graveyardCase := range analysis.Graveyard.Cases {
			report.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, graveyardCase.CompanyName))
			report.WriteString(fmt.Sprintf("   - **Description:** %s\n", graveyardCase.Description))
			report.WriteString(fmt.Sprintf("   - **Failure Cause:** %s\n", graveyardCase.FailureCause))
			report.WriteString(fmt.Sprintf("   - **Lessons:** %s\n", graveyardCase.Lessons))
			if len(graveyardCase.EvidenceIDs) > 0 {
				report.WriteString(fmt.Sprintf("   - Sources: %s\n", mb.formatEvidenceRefs(graveyardCase.EvidenceIDs)))
			}
			report.WriteString("\n")
		}
	}

	// Evidence References
	if len(analysis.Evidence) > 0 {
		report.WriteString("## Evidence References\n\n")
		evidenceMap := make(map[string]types.Evidence)
		for _, ev := range analysis.Evidence {
			evidenceMap[ev.ID] = ev
		}

		counter := 1
		for _, ev := range analysis.Evidence {
			report.WriteString(fmt.Sprintf("[%d] **%s**\n", counter, ev.Title))
			report.WriteString(fmt.Sprintf("    %s\n", ev.URL))
			if ev.Snippet != "" {
				report.WriteString(fmt.Sprintf("    %s\n", ev.Snippet))
			}
			if ev.PublishedAt != nil {
				report.WriteString(fmt.Sprintf("    Published: %s\n", ev.PublishedAt.Format("January 2, 2006")))
			}
			report.WriteString(fmt.Sprintf("    Source: %s\n", strings.Title(ev.SourceType)))
			report.WriteString("\n")
			counter++
		}
	}

	// Footer
	report.WriteString("---\n\n")
	report.WriteString("*Generated by RectAIfy*\n")

	return report.String()
}

// getScoreAssessment returns a textual assessment based on score
func (mb *MarkdownBuilder) getScoreAssessment(score float64) string {
	if score >= 80 {
		return "Excellent"
	} else if score >= 60 {
		return "Good"
	} else if score >= 40 {
		return "Fair"
	} else if score >= 20 {
		return "Poor"
	} else {
		return "Critical"
	}
}

// formatEvidenceRefs formats evidence IDs as numbered references
func (mb *MarkdownBuilder) formatEvidenceRefs(evidenceIDs []string) string {
	if len(evidenceIDs) == 0 {
		return ""
	}

	refs := make([]string, len(evidenceIDs))
	for i := range evidenceIDs {
		// For now, just use index+1. In a full implementation,
		// we'd maintain a mapping from evidence ID to reference number
		refs[i] = fmt.Sprintf("[%d]", i+1)
	}

	return strings.Join(refs, ", ")
}
