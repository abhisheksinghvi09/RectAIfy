package score

import (
	"math"
	"strings"

	"rectaify/pkg/types"
)

// Calculator computes viability scores based on analysis results
type Calculator struct {
	weights ScoreWeights
}

// ScoreWeights defines the relative importance of each scoring dimension
type ScoreWeights struct {
	Market     float64 `json:"market"`
	Problem    float64 `json:"problem"`
	Barriers   float64 `json:"barriers"`
	Execution  float64 `json:"execution"`
	Risks      float64 `json:"risks"`
	Graveyard  float64 `json:"graveyard"`
}

// DefaultWeights returns sensible default weights
func DefaultWeights() ScoreWeights {
	return ScoreWeights{
		Market:    0.25, // 25% - Market opportunity and competition
		Problem:   0.20, // 20% - Problem validation
		Barriers:  0.15, // 15% - Execution barriers
		Execution: 0.15, // 15% - Execution complexity
		Risks:     0.15, // 15% - Business risks
		Graveyard: 0.10, // 10% - Learning from failures
	}
}

// NewCalculator creates a new score calculator
func NewCalculator(weights *ScoreWeights) *Calculator {
	if weights == nil {
		defaultWeights := DefaultWeights()
		weights = &defaultWeights
	}
	return &Calculator{weights: *weights}
}

// ComputeViability calculates the overall viability score
func (c *Calculator) ComputeViability(analysis types.Analysis) types.Viability {
	marketScore := c.computeMarketScore(analysis.Market)
	problemScore := c.computeProblemScore(analysis.Problem)
	barrierScore := c.computeBarrierScore(analysis.Barriers)
	executionScore := c.computeExecutionScore(analysis.Execution)
	riskScore := c.computeRiskScore(analysis.Risks)
	graveyardScore := c.computeGraveyardScore(analysis.Graveyard)

	// Calculate weighted overall score
	overallScore := (marketScore * c.weights.Market) +
		(problemScore * c.weights.Problem) +
		(barrierScore * c.weights.Barriers) +
		(executionScore * c.weights.Execution) +
		(riskScore * c.weights.Risks) +
		(graveyardScore * c.weights.Graveyard)

	// Ensure score is bounded [0, 100]
	overallScore = math.Max(0, math.Min(100, overallScore))

	recommendation := c.generateRecommendation(overallScore, marketScore, problemScore, barrierScore, executionScore, riskScore, graveyardScore)
	keyInsights := c.generateKeyInsights(analysis, marketScore, problemScore, barrierScore, executionScore, riskScore, graveyardScore)

	// Collect all evidence IDs
	evidenceIDs := c.collectEvidenceIDs(analysis)

	return types.Viability{
		OverallScore:    overallScore,
		MarketScore:     marketScore,
		ProblemScore:    problemScore,
		BarrierScore:    barrierScore,
		ExecutionScore:  executionScore,
		RiskScore:       riskScore,
		GraveyardScore:  graveyardScore,
		Recommendation:  recommendation,
		KeyInsights:     keyInsights,
		EvidenceIDs:     evidenceIDs,
	}
}

// computeMarketScore calculates market opportunity score
func (c *Calculator) computeMarketScore(market types.MarketAnalysis) float64 {
	score := 50.0 // Base score

	// Stage scoring
	stageScores := map[string]float64{
		"early":     85.0, // High opportunity in early markets
		"growing":   70.0, // Good opportunity in growing markets
		"mature":    40.0, // Harder in mature markets
		"declining": 15.0, // Very difficult in declining markets
	}

	if stageScore, exists := stageScores[market.MarketStage]; exists {
		score = stageScore
	}

	// Competition adjustment
	competitorCount := len(market.Competitors)
	if competitorCount == 0 {
		score += 15.0 // Blue ocean opportunity
	} else if competitorCount <= 2 {
		score += 5.0 // Limited competition
	} else if competitorCount <= 5 {
		score -= 5.0 // Moderate competition
	} else {
		score -= 15.0 // High competition
	}

	// Positioning quality
	if market.Positioning != "" {
		if len(market.Positioning) > 50 {
			score += 5.0 // Well-defined positioning
		}
	}

	// Evidence quality bonus
	evidenceBonus := math.Min(10.0, float64(len(market.EvidenceIDs))*2.0)
	score += evidenceBonus

	return math.Max(0, math.Min(100, score))
}

// computeProblemScore calculates problem validation score
func (c *Calculator) computeProblemScore(problem types.ProblemAnalysis) float64 {
	score := 30.0 // Base score (problems need validation)

	// Pain points count
	painPointCount := len(problem.PainPoints)
	if painPointCount >= 3 {
		score += 25.0 // Multiple clear pain points
	} else if painPointCount >= 2 {
		score += 15.0 // Some pain points
	} else if painPointCount >= 1 {
		score += 10.0 // At least one pain point
	}

	// Validation quality
	if problem.Validation != "" {
		validationLength := len(problem.Validation)
		if validationLength > 100 {
			score += 20.0 // Strong validation
		} else if validationLength > 50 {
			score += 10.0 // Some validation
		}
	}

	// Evidence quality bonus
	evidenceBonus := math.Min(15.0, float64(len(problem.EvidenceIDs))*3.0)
	score += evidenceBonus

	return math.Max(0, math.Min(100, score))
}

// computeBarrierScore calculates execution barrier score (lower barriers = higher score)
func (c *Calculator) computeBarrierScore(barriers types.BarrierAnalysis) float64 {
	if len(barriers.Barriers) == 0 {
		return 85.0 // No significant barriers identified
	}

	// Calculate weighted barrier impact
	totalWeight := 0.0
	weightedImpact := 0.0

	for _, barrier := range barriers.Barriers {
		totalWeight += barrier.Weight
		
		// Convert barrier type to impact score
		barrierImpact := c.getBarrierImpact(barrier.Type)
		weightedImpact += barrier.Weight * barrierImpact
	}

	if totalWeight == 0 {
		return 85.0
	}

	// Average weighted impact (0-100, where 100 is highest barrier)
	avgImpact := weightedImpact / totalWeight

	// Convert to score (inverse relationship - lower barriers = higher score)
	score := 100.0 - avgImpact

	// Evidence adjustment
	evidenceCount := len(barriers.EvidenceIDs)
	if evidenceCount > 0 {
		// More evidence of barriers = more reliable assessment
		reliabilityBonus := math.Min(5.0, float64(evidenceCount))
		score -= reliabilityBonus // Subtract because more evidence of barriers is bad
	}

	return math.Max(0, math.Min(100, score))
}

// getBarrierImpact returns impact score for different barrier types
func (c *Calculator) getBarrierImpact(barrierType string) float64 {
	impacts := map[string]float64{
		"regulation":   85.0, // Very high impact
		"supply":       70.0, // High impact
		"distribution": 60.0, // Moderate-high impact
		"trust":        50.0, // Moderate impact
		"tech":         40.0, // Moderate-low impact
	}

	if impact, exists := impacts[barrierType]; exists {
		return impact
	}
	return 50.0 // Default moderate impact
}

// computeExecutionScore calculates execution complexity score
func (c *Calculator) computeExecutionScore(execution types.ExecutionAnalysis) float64 {
	score := 70.0 // Base score

	// Capital requirement impact
	capitalScores := map[string]float64{
		"low":    90.0,
		"medium": 60.0,
		"high":   30.0,
		"very high": 10.0,
	}

	if capitalScore, exists := capitalScores[execution.CapitalRequirement]; exists {
		score = (score + capitalScore) / 2.0
	}

	// Talent rarity impact
	talentScores := map[string]float64{
		"common":    85.0,
		"available": 70.0,
		"scarce":    45.0,
		"rare":      25.0,
	}

	if talentScore, exists := talentScores[execution.TalentRarity]; exists {
		score = (score + talentScore) / 2.0
	}

	// Integration complexity (more integrations = lower score)
	integrationPenalty := math.Min(30.0, float64(execution.IntegrationCount)*5.0)
	score -= integrationPenalty

	// Direct complexity score
	if execution.Complexity > 0 {
		complexityScore := 100.0 - (execution.Complexity * 100.0)
		score = (score + complexityScore) / 2.0
	}

	// Evidence quality adjustment
	evidenceBonus := math.Min(5.0, float64(len(execution.EvidenceIDs)))
	score += evidenceBonus

	return math.Max(0, math.Min(100, score))
}

// computeRiskScore calculates business risk score
func (c *Calculator) computeRiskScore(risks types.RiskAnalysis) float64 {
	if len(risks.Risks) == 0 {
		return 80.0 // No identified risks (but this might be bad research)
	}

	score := 100.0 // Start high, subtract for risks

	totalRiskImpact := 0.0
	riskCount := 0

	for _, risk := range risks.Risks {
		// Calculate risk impact (severity * likelihood)
		impact := float64(risk.Severity * risk.Likelihood) // Max is 25 (5*5)
		totalRiskImpact += impact
		riskCount++

		// Deduct based on risk impact
		riskPenalty := (impact / 25.0) * 20.0 // Scale to max 20 points per risk
		score -= riskPenalty

		// Mitigation bonus
		if risk.Mitigation != "" && len(risk.Mitigation) > 20 {
			score += 3.0 // Small bonus for having mitigation plans
		}
	}

	// Evidence quality adjustment
	evidenceCount := len(risks.EvidenceIDs)
	if evidenceCount > 0 {
		reliabilityBonus := math.Min(5.0, float64(evidenceCount))
		score += reliabilityBonus
	}

	return math.Max(0, math.Min(100, score))
}

// computeGraveyardScore calculates learning from failures score
func (c *Calculator) computeGraveyardScore(graveyard types.GraveyardAnalysis) float64 {
	if len(graveyard.Cases) == 0 {
		return 60.0 // No failure cases found - could be good or bad
	}

	score := 40.0 // Start lower when failures exist

	for _, graveyardCase := range graveyard.Cases {
		// Penalty for each failure case
		score -= 10.0

		// Bonus for having lessons learned
		if graveyardCase.Lessons != "" && len(graveyardCase.Lessons) > 30 {
			score += 5.0 // Learning from failures is valuable
		}

		// Check failure cause patterns
		cause := strings.ToLower(graveyardCase.FailureCause)
		if strings.Contains(cause, "funding") || strings.Contains(cause, "money") {
			score -= 5.0 // Funding failures are concerning
		} else if strings.Contains(cause, "market") || strings.Contains(cause, "demand") {
			score -= 8.0 // Market failures are very concerning
		} else if strings.Contains(cause, "execution") || strings.Contains(cause, "team") {
			score -= 3.0 // Execution failures are somewhat concerning
		}
	}

	// Evidence quality bonus
	evidenceBonus := math.Min(10.0, float64(len(graveyard.EvidenceIDs))*2.0)
	score += evidenceBonus

	return math.Max(0, math.Min(100, score))
}

// generateRecommendation creates a recommendation based on scores
func (c *Calculator) generateRecommendation(overall, market, problem, barrier, execution, risk, graveyard float64) string {
	if overall >= 75 {
		return "STRONG GO: High viability with favorable conditions across multiple dimensions."
	} else if overall >= 60 {
		return "GO: Good viability with some areas requiring attention."
	} else if overall >= 45 {
		return "CAUTION: Mixed signals - proceed with careful validation and risk mitigation."
	} else if overall >= 30 {
		return "HIGH RISK: Significant challenges identified - major pivots likely needed."
	} else {
		return "NO GO: Multiple severe challenges make success highly unlikely."
	}
}

// generateKeyInsights extracts key insights from the scoring analysis
func (c *Calculator) generateKeyInsights(analysis types.Analysis, market, problem, barrier, execution, risk, graveyard float64) []string {
	var insights []string

	// Market insights
	if market >= 80 {
		insights = append(insights, "Strong market opportunity with favorable competitive dynamics")
	} else if market <= 30 {
		insights = append(insights, "Challenging market conditions with intense competition or declining demand")
	}

	// Problem insights
	if problem >= 80 {
		insights = append(insights, "Well-validated problem with clear pain points")
	} else if problem <= 40 {
		insights = append(insights, "Problem validation is weak - more research needed")
	}

	// Barrier insights
	if barrier <= 40 {
		insights = append(insights, "Significant execution barriers identified")
	} else if barrier >= 80 {
		insights = append(insights, "Clear path to execution with minimal barriers")
	}

	// Execution insights
	if execution <= 40 {
		insights = append(insights, "High execution complexity requiring substantial resources")
	} else if execution >= 80 {
		insights = append(insights, "Manageable execution complexity with available resources")
	}

	// Risk insights
	if risk <= 40 {
		insights = append(insights, "High business risks requiring careful mitigation")
	} else if risk >= 80 {
		insights = append(insights, "Well-managed risk profile")
	}

	// Graveyard insights
	if graveyard <= 40 && len(analysis.Graveyard.Cases) > 0 {
		insights = append(insights, "Multiple similar ventures have failed - learn from their mistakes")
	}

	// Ensure we have at least one insight
	if len(insights) == 0 {
		insights = append(insights, "Further research recommended to validate assumptions")
	}

	return insights
}

// collectEvidenceIDs gathers all evidence IDs from the analysis
func (c *Calculator) collectEvidenceIDs(analysis types.Analysis) []string {
	evidenceMap := make(map[string]bool)
	
	// Collect from all analysis sections
	for _, id := range analysis.Market.EvidenceIDs {
		evidenceMap[id] = true
	}
	for _, id := range analysis.Problem.EvidenceIDs {
		evidenceMap[id] = true
	}
	for _, id := range analysis.Barriers.EvidenceIDs {
		evidenceMap[id] = true
	}
	for _, id := range analysis.Execution.EvidenceIDs {
		evidenceMap[id] = true
	}
	for _, id := range analysis.Risks.EvidenceIDs {
		evidenceMap[id] = true
	}
	for _, id := range analysis.Graveyard.EvidenceIDs {
		evidenceMap[id] = true
	}

	// Convert to slice
	var evidenceIDs []string
	for id := range evidenceMap {
		evidenceIDs = append(evidenceIDs, id)
	}

	return evidenceIDs
}
