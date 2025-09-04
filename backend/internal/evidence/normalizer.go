package evidence

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"realitycheck/pkg/types"
)

// Normalizer handles evidence normalization and deduplication
type Normalizer struct {
	minHashSize int
}

// NewNormalizer creates a new evidence normalizer
func NewNormalizer() *Normalizer {
	return &Normalizer{
		minHashSize: 3, // MinHash signature size
	}
}

// Normalize processes and normalizes evidence
func (n *Normalizer) Normalize(evidence []types.Evidence) []types.Evidence {
	// First pass: normalize individual evidence entries
	normalized := make([]types.Evidence, 0, len(evidence))
	for _, ev := range evidence {
		if normalizedEv := n.normalizeEvidence(ev); normalizedEv != nil {
			normalized = append(normalized, *normalizedEv)
		}
	}

	// Second pass: deduplicate similar evidence
	deduped := n.deduplicateEvidence(normalized)

	// Third pass: quality filtering and ranking
	filtered := n.filterByQuality(deduped)

	return filtered
}

// normalizeEvidence normalizes a single evidence entry
func (n *Normalizer) normalizeEvidence(ev types.Evidence) *types.Evidence {
	// Validate required fields
	if ev.URL == "" || ev.Title == "" {
		return nil
	}

	// Canonicalize URL
	canonicalURL := n.canonicalizeURL(ev.URL)
	if canonicalURL == "" {
		return nil // Invalid URL
	}

	// Clean title and snippet
	cleanTitle := n.cleanText(ev.Title)
	cleanSnippet := n.cleanText(ev.Snippet)

	// Generate stable ID
	stableID := n.generateStableID(canonicalURL, cleanTitle, ev.PublishedAt)

	// Infer source type if not provided
	sourceType := ev.SourceType
	if sourceType == "" {
		sourceType = n.inferSourceType(canonicalURL)
	}

	return &types.Evidence{
		ID:          stableID,
		URL:         canonicalURL,
		Title:       cleanTitle,
		Snippet:     cleanSnippet,
		PublishedAt: ev.PublishedAt,
		RetrievedAt: ev.RetrievedAt,
		SourceType:  sourceType,
	}
}

// canonicalizeURL normalizes URLs by removing tracking parameters
func (n *Normalizer) canonicalizeURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	// Reject invalid schemes
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}

	// Remove common tracking parameters
	trackingParams := []string{
		"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content",
		"gclid", "fbclid", "msclkid", "ref", "referrer", "source",
		"_ga", "_gl", "mc_cid", "mc_eid", "WT.mc_id",
		"campaign", "medium", "content", "term",
	}

	q := u.Query()
	for _, param := range trackingParams {
		q.Del(param)
	}
	u.RawQuery = q.Encode()

	// Normalize host (remove www.)
	if strings.HasPrefix(u.Host, "www.") {
		u.Host = u.Host[4:]
	}

	return u.String()
}

// cleanText cleans and normalizes text content
func (n *Normalizer) cleanText(text string) string {
	if text == "" {
		return ""
	}

	// Remove excessive whitespace
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Normalize multiple spaces
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	// Limit length
	maxLength := 500
	if len(text) > maxLength {
		text = text[:maxLength] + "..."
	}

	return text
}

// generateStableID creates a stable ID for evidence
func (n *Normalizer) generateStableID(url, title string, publishedAt *time.Time) string {
	var timeStr string
	if publishedAt != nil {
		timeStr = publishedAt.Format(time.RFC3339)
	}

	content := fmt.Sprintf("%s|%s|%s", url, title, timeStr)
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes for shorter ID
}

// inferSourceType determines the source type from URL
func (n *Normalizer) inferSourceType(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "unknown"
	}

	domain := strings.ToLower(u.Host)

	// Map domains to source types
	sourceTypes := map[string]string{
		"techcrunch.com":      "news",
		"venturebeat.com":     "news",
		"arstechnica.com":     "news",
		"theverge.com":        "news",
		"wired.com":           "news",
		"reuters.com":         "news",
		"bloomberg.com":       "news",
		"wsj.com":             "news",
		"nytimes.com":         "news",
		"forbes.com":          "news",
		"fortune.com":         "news",
		"businessinsider.com": "news",
		"crunchbase.com":      "database",
		"pitchbook.com":       "database",
		"sec.gov":             "regulatory",
		"fda.gov":             "regulatory",
		"reddit.com":          "forum",
		"news.ycombinator.com": "forum",
		"github.com":          "code",
		"stackoverflow.com":   "forum",
		"medium.com":          "blog",
		"substack.com":        "blog",
		"linkedin.com":        "professional",
		"twitter.com":         "social",
		"x.com":               "social",
		"youtube.com":         "video",
		"angellist.com":       "startup",
		"wellfound.com":       "startup",
		"producthunt.com":     "product",
		"ycombinator.com":     "accelerator",
		"techstars.com":       "accelerator",
	}

	if sourceType, exists := sourceTypes[domain]; exists {
		return sourceType
	}

	// Default categorization based on patterns
	if strings.Contains(domain, "gov") {
		return "government"
	}
	if strings.Contains(domain, "edu") {
		return "academic"
	}
	if strings.Contains(domain, "blog") {
		return "blog"
	}
	if strings.Contains(domain, "news") {
		return "news"
	}

	return "website"
}

// deduplicateEvidence removes near-duplicate evidence using multiple strategies
func (n *Normalizer) deduplicateEvidence(evidence []types.Evidence) []types.Evidence {
	if len(evidence) <= 1 {
		return evidence
	}

	// Group by URL+title first (exact duplicates)
	urlTitleMap := make(map[string]types.Evidence)
	for _, ev := range evidence {
		key := ev.URL + "|" + ev.Title
		if existing, exists := urlTitleMap[key]; exists {
			// Keep the one with more recent publication date
			if ev.PublishedAt != nil && (existing.PublishedAt == nil || ev.PublishedAt.After(*existing.PublishedAt)) {
				urlTitleMap[key] = ev
			}
		} else {
			urlTitleMap[key] = ev
		}
	}

	// Convert back to slice
	unique := make([]types.Evidence, 0, len(urlTitleMap))
	for _, ev := range urlTitleMap {
		unique = append(unique, ev)
	}

	// Apply content similarity deduplication
	filtered := n.filterSimilarContent(unique)

	return filtered
}

// filterSimilarContent removes evidence with very similar content
func (n *Normalizer) filterSimilarContent(evidence []types.Evidence) []types.Evidence {
	if len(evidence) <= 1 {
		return evidence
	}

	var filtered []types.Evidence
	processed := make(map[int]bool)

	for i, ev1 := range evidence {
		if processed[i] {
			continue
		}

		// Find all similar evidence
		similar := []int{i}
		for j := i + 1; j < len(evidence); j++ {
			if processed[j] {
				continue
			}

			ev2 := evidence[j]
			if n.areContentSimilar(ev1, ev2) {
				similar = append(similar, j)
			}
		}

		// Mark all as processed
		for _, idx := range similar {
			processed[idx] = true
		}

		// Select the best representative from similar group
		best := n.selectBestEvidence(evidence, similar)
		filtered = append(filtered, best)
	}

	return filtered
}

// areContentSimilar determines if two evidence entries have similar content
func (n *Normalizer) areContentSimilar(ev1, ev2 types.Evidence) bool {
	// Check title similarity
	titleSim := n.textSimilarity(ev1.Title, ev2.Title)
	if titleSim > 0.8 {
		return true
	}

	// Check snippet similarity if both have snippets
	if ev1.Snippet != "" && ev2.Snippet != "" {
		snippetSim := n.textSimilarity(ev1.Snippet, ev2.Snippet)
		if snippetSim > 0.7 {
			return true
		}
	}

	// Check if they're from the same domain with similar titles
	domain1 := n.extractDomain(ev1.URL)
	domain2 := n.extractDomain(ev2.URL)
	if domain1 == domain2 && titleSim > 0.6 {
		return true
	}

	return false
}

// textSimilarity calculates simple text similarity using Jaccard index
func (n *Normalizer) textSimilarity(text1, text2 string) float64 {
	if text1 == "" || text2 == "" {
		return 0
	}

	words1 := n.tokenize(text1)
	words2 := n.tokenize(text2)

	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, word := range words1 {
		set1[word] = true
	}

	for _, word := range words2 {
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

// tokenize splits text into normalized tokens
func (n *Normalizer) tokenize(text string) []string {
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'))
	})

	// Filter out short words and common stop words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true,
	}

	var filtered []string
	for _, word := range words {
		if len(word) > 2 && !stopWords[word] {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// extractDomain extracts domain from URL
func (n *Normalizer) extractDomain(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return u.Host
}

// selectBestEvidence selects the best evidence from a group of similar ones
func (n *Normalizer) selectBestEvidence(evidence []types.Evidence, indices []int) types.Evidence {
	if len(indices) == 1 {
		return evidence[indices[0]]
	}

	// Score each evidence based on quality factors
	best := evidence[indices[0]]
	bestScore := n.scoreEvidenceQuality(best)

	for i := 1; i < len(indices); i++ {
		ev := evidence[indices[i]]
		score := n.scoreEvidenceQuality(ev)
		if score > bestScore {
			best = ev
			bestScore = score
		}
	}

	return best
}

// scoreEvidenceQuality assigns a quality score to evidence
func (n *Normalizer) scoreEvidenceQuality(ev types.Evidence) float64 {
	score := 0.0

	// Source type scoring
	sourceScores := map[string]float64{
		"news":        1.0,
		"database":    0.9,
		"regulatory":  0.9,
		"academic":    0.8,
		"professional": 0.7,
		"startup":     0.7,
		"code":        0.6,
		"blog":        0.5,
		"forum":       0.4,
		"social":      0.3,
		"video":       0.3,
		"website":     0.2,
		"unknown":     0.1,
	}

	if sourceScore, exists := sourceScores[ev.SourceType]; exists {
		score += sourceScore
	}

	// Published date scoring (more recent = better)
	if ev.PublishedAt != nil {
		daysSince := time.Since(*ev.PublishedAt).Hours() / 24
		if daysSince <= 30 {
			score += 0.5 // Very recent
		} else if daysSince <= 365 {
			score += 0.3 // Recent
		} else if daysSince <= 365*3 {
			score += 0.1 // Somewhat recent
		}
	}

	// Content quality scoring
	if len(ev.Title) > 10 {
		score += 0.2
	}
	if len(ev.Snippet) > 50 {
		score += 0.2
	}

	// URL quality (shorter is often better)
	if len(ev.URL) < 100 {
		score += 0.1
	}

	return score
}

// filterByQuality removes low-quality evidence and sorts by quality
func (n *Normalizer) filterByQuality(evidence []types.Evidence) []types.Evidence {
	// Score all evidence
	type scoredEvidence struct {
		evidence types.Evidence
		score    float64
	}

	scored := make([]scoredEvidence, 0, len(evidence))
	for _, ev := range evidence {
		score := n.scoreEvidenceQuality(ev)
		if score > 0.3 { // Minimum quality threshold
			scored = append(scored, scoredEvidence{evidence: ev, score: score})
		}
	}

	// Sort by score (highest first)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Extract evidence
	filtered := make([]types.Evidence, len(scored))
	for i, se := range scored {
		filtered[i] = se.evidence
	}

	return filtered
}
