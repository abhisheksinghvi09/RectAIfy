package llm

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// generateEvidenceID creates a stable ID for evidence based on content
func generateEvidenceID(urlStr, title string, publishedAt *time.Time) string {
	var timeStr string
	if publishedAt != nil {
		timeStr = publishedAt.Format(time.RFC3339)
	}
	
	content := fmt.Sprintf("%s|%s|%s", urlStr, title, timeStr)
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes for shorter ID
}

// inferSourceType determines the type of source based on URL
func inferSourceType(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "unknown"
	}

	domain := strings.ToLower(u.Host)
	
	// Remove www. prefix
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:]
	}

	// Map domains to source types
	sourceTypes := map[string]string{
		"techcrunch.com":      "news",
		"venturebeat.com":     "news",
		"techstars.com":       "accelerator",
		"ycombinator.com":     "accelerator",
		"crunchbase.com":      "database",
		"pitchbook.com":       "database",
		"sec.gov":             "regulatory",
		"reddit.com":          "forum",
		"hackernews.com":      "forum",
		"github.com":          "code",
		"medium.com":          "blog",
		"substack.com":        "blog",
		"linkedin.com":        "professional",
		"twitter.com":         "social",
		"x.com":               "social",
		"youtube.com":         "video",
		"angellist.com":       "startup",
		"wellfound.com":       "startup",
		"producthunt.com":     "product",
		"reuters.com":         "news",
		"bloomberg.com":       "news",
		"wsj.com":             "news",
		"nytimes.com":         "news",
		"washingtonpost.com":  "news",
		"forbes.com":          "news",
		"fortune.com":         "news",
		"businessinsider.com": "news",
		"wired.com":           "news",
		"arstechnica.com":     "news",
		"theverge.com":        "news",
	}

	if sourceType, exists := sourceTypes[domain]; exists {
		return sourceType
	}

	// Default categorization based on TLD or patterns
	if strings.Contains(domain, "gov") {
		return "government"
	}
	if strings.Contains(domain, "edu") {
		return "academic"
	}
	if strings.Contains(domain, "blog") || strings.Contains(domain, "medium") {
		return "blog"
	}
	if strings.Contains(domain, "news") {
		return "news"
	}

	return "website"
}

// canonicalizeURL normalizes URLs by removing tracking parameters
func canonicalizeURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	// Remove common tracking parameters
	trackingParams := []string{
		"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content",
		"gclid", "fbclid", "msclkid", "ref", "referrer",
		"_ga", "_gl", "mc_cid", "mc_eid",
	}

	q := u.Query()
	for _, param := range trackingParams {
		q.Del(param)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
