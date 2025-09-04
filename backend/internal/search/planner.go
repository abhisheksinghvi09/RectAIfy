package search

import (
	"context"
	"fmt"
	"strings"

	"rectaify/pkg/types"
)

// Planner generates search queries from startup ideas
type Planner struct {
	maxQueries int
}

// NewPlanner creates a new query planner
func NewPlanner(maxQueries int) *Planner {
	return &Planner{
		maxQueries: maxQueries,
	}
}

// Plan generates search queries from an idea
func (p *Planner) Plan(ctx context.Context, idea types.IdeaInput) ([]types.SearchQuery, error) {
	var queries []types.SearchQuery
	
	// Normalize the idea text
	normalizedTitle := normalizeText(idea.Title)
	normalizedOneLiner := normalizeText(idea.OneLiner)
	
	// Extract key terms
	keyTerms := extractKeyTerms(normalizedTitle, normalizedOneLiner)
	
	// Generate queries by intent
	queries = append(queries, p.generateCompetitorQueries(keyTerms, idea)...)
	queries = append(queries, p.generateFundingQueries(keyTerms, idea)...)
	queries = append(queries, p.generateRegulatoryQueries(keyTerms, idea)...)
	queries = append(queries, p.generatePostmortemQueries(keyTerms, idea)...)
	queries = append(queries, p.generateMarketQueries(keyTerms, idea)...)
	queries = append(queries, p.generateProblemQueries(keyTerms, idea)...)
	
	// Deduplicate and limit
	queries = p.deduplicateQueries(queries)
	
	if len(queries) > p.maxQueries {
		queries = queries[:p.maxQueries]
	}
	
	return queries, nil
}

// generateCompetitorQueries creates queries to find competitors
func (p *Planner) generateCompetitorQueries(keyTerms []string, idea types.IdeaInput) []types.SearchQuery {
	var queries []types.SearchQuery
	
	templates := []string{
		"%s competitors",
		"%s alternative",
		"%s similar companies",
		"companies like %s",
		"%s vs competitors",
		"best %s tools",
		"%s market leaders",
		"top %s startups",
	}
	
	for _, term := range keyTerms[:min(len(keyTerms), 3)] {
		for _, template := range templates[:4] { // Limit templates
			query := fmt.Sprintf(template, term)
			queries = append(queries, types.SearchQuery{
				Query:    query,
				Intent:   "competitors",
				Priority: 1,
			})
		}
	}
	
	// Add specific queries based on the idea
	queries = append(queries, types.SearchQuery{
		Query:    fmt.Sprintf("\"%s\" competitors", idea.Title),
		Intent:   "competitors",
		Priority: 2,
	})
	
	return queries
}

// generateFundingQueries creates queries to find funding information
func (p *Planner) generateFundingQueries(keyTerms []string, idea types.IdeaInput) []types.SearchQuery {
	var queries []types.SearchQuery
	
	templates := []string{
		"%s startup funding",
		"%s series A",
		"%s investment",
		"%s venture capital",
		"%s raised money",
		"funding %s startups",
		"%s IPO",
		"%s acquisition",
	}
	
	for _, term := range keyTerms[:min(len(keyTerms), 2)] {
		for _, template := range templates[:4] {
			query := fmt.Sprintf(template, term)
			queries = append(queries, types.SearchQuery{
				Query:    query,
				Intent:   "funding",
				Priority: 2,
			})
		}
	}
	
	return queries
}

// generateRegulatoryQueries creates queries to find regulatory information
func (p *Planner) generateRegulatoryQueries(keyTerms []string, idea types.IdeaInput) []types.SearchQuery {
	var queries []types.SearchQuery
	
	templates := []string{
		"%s regulation",
		"%s compliance",
		"%s legal requirements",
		"%s government rules",
		"%s licensing",
		"%s permits",
		"%s regulatory approval",
		"%s FDA approval",
	}
	
	for _, term := range keyTerms[:min(len(keyTerms), 2)] {
		for _, template := range templates[:4] {
			query := fmt.Sprintf(template, term)
			queries = append(queries, types.SearchQuery{
				Query:    query,
				Intent:   "regulation",
				Priority: 2,
			})
		}
	}
	
	return queries
}

// generatePostmortemQueries creates queries to find failure cases
func (p *Planner) generatePostmortemQueries(keyTerms []string, idea types.IdeaInput) []types.SearchQuery {
	var queries []types.SearchQuery
	
	templates := []string{
		"%s startup failed",
		"%s company shut down",
		"%s startup postmortem",
		"why %s failed",
		"%s startup lessons",
		"failed %s companies",
		"%s startup mistakes",
		"%s business failed",
	}
	
	for _, term := range keyTerms[:min(len(keyTerms), 2)] {
		for _, template := range templates[:4] {
			query := fmt.Sprintf(template, term)
			queries = append(queries, types.SearchQuery{
				Query:    query,
				Intent:   "postmortems",
				Priority: 3,
			})
		}
	}
	
	return queries
}

// generateMarketQueries creates queries to understand market size and trends
func (p *Planner) generateMarketQueries(keyTerms []string, idea types.IdeaInput) []types.SearchQuery {
	var queries []types.SearchQuery
	
	templates := []string{
		"%s market size",
		"%s industry trends",
		"%s market research",
		"%s TAM",
		"%s market opportunity",
		"global %s market",
		"%s industry analysis",
		"%s market growth",
	}
	
	for _, term := range keyTerms[:min(len(keyTerms), 2)] {
		for _, template := range templates[:4] {
			query := fmt.Sprintf(template, term)
			queries = append(queries, types.SearchQuery{
				Query:    query,
				Intent:   "market",
				Priority: 1,
			})
		}
	}
	
	return queries
}

// generateProblemQueries creates queries to validate the problem
func (p *Planner) generateProblemQueries(keyTerms []string, idea types.IdeaInput) []types.SearchQuery {
	var queries []types.SearchQuery
	
	templates := []string{
		"%s problems",
		"%s pain points",
		"users complain %s",
		"%s frustrations",
		"why %s sucks",
		"%s issues",
		"problems with %s",
		"%s challenges",
	}
	
	for _, term := range keyTerms[:min(len(keyTerms), 2)] {
		for _, template := range templates[:4] {
			query := fmt.Sprintf(template, term)
			queries = append(queries, types.SearchQuery{
				Query:    query,
				Intent:   "problem",
				Priority: 1,
			})
		}
	}
	
	return queries
}

// deduplicateQueries removes similar queries using token set similarity
func (p *Planner) deduplicateQueries(queries []types.SearchQuery) []types.SearchQuery {
	if len(queries) <= 1 {
		return queries
	}
	
	var unique []types.SearchQuery
	seen := make(map[string]bool)
	
	for _, query := range queries {
		// Normalize query for comparison
		normalized := normalizeQuery(query.Query)
		
		// Check for duplicates
		isDuplicate := false
		for existing := range seen {
			if jaccardSimilarity(normalized, existing) > 0.8 {
				isDuplicate = true
				break
			}
		}
		
		if !isDuplicate {
			seen[normalized] = true
			unique = append(unique, query)
		}
	}
	
	return unique
}

// normalizeText cleans and normalizes text
func normalizeText(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	
	// Remove common stop words and punctuation
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "is": true,
		"are": true, "was": true, "were": true, "be": true, "been": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
	}
	
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'))
	})
	
	var filtered []string
	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			filtered = append(filtered, word)
		}
	}
	
	return strings.Join(filtered, " ")
}

// extractKeyTerms extracts important terms from the idea
func extractKeyTerms(title, oneLiner string) []string {
	allText := title + " " + oneLiner
	words := strings.Fields(allText)
	
	// Simple term extraction - take longer words and capitalize words
	var keyTerms []string
	seen := make(map[string]bool)
	
	for _, word := range words {
		word = strings.ToLower(word)
		
		// Skip if already seen, too short, or common
		if seen[word] || len(word) < 3 {
			continue
		}
		
		// Add significant terms
		if len(word) >= 5 || strings.Title(word) == word {
			keyTerms = append(keyTerms, word)
			seen[word] = true
		}
	}
	
	return keyTerms
}

// normalizeQuery normalizes a query for comparison
func normalizeQuery(query string) string {
	// Convert to lowercase and extract words
	words := strings.FieldsFunc(strings.ToLower(query), func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'))
	})
	
	// Sort words for consistent comparison
	return strings.Join(words, " ")
}

// jaccardSimilarity calculates Jaccard similarity between two queries
func jaccardSimilarity(query1, query2 string) float64 {
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, word := range strings.Fields(query1) {
		set1[word] = true
	}
	
	for _, word := range strings.Fields(query2) {
		set2[word] = true
	}
	
	intersection := 0
	for word := range set1 {
		if set2[word] {
			intersection++
		}
	}
	
	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
