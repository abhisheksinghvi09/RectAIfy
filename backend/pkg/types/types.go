package types

import (
	"encoding/json"
	"time"
)

// IdeaInput represents the initial startup idea to be analyzed
type IdeaInput struct {
	Title    string `json:"title" validate:"required,min=1,max=200"`
	OneLiner string `json:"one_liner" validate:"required,min=10,max=500"`
	Category string `json:"category,omitempty"`
	Location string `json:"location,omitempty"` // for geographic context
}

// Evidence represents a piece of research evidence with citations
type Evidence struct {
	ID          string     `json:"id" db:"id"`
	URL         string     `json:"url" db:"url"`
	Title       string     `json:"title" db:"title"`
	Snippet     string     `json:"snippet,omitempty" db:"snippet"`
	PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at"`
	RetrievedAt time.Time  `json:"retrieved_at" db:"retrieved_at"`
	SourceType  string     `json:"source_type,omitempty" db:"source_type"`
}

// Competitor represents market competition analysis
type Competitor struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Funding     string   `json:"funding,omitempty"`
	Stage       string   `json:"stage,omitempty"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// Risk represents identified business risks
type Risk struct {
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Severity    int      `json:"severity"` // 1-5 scale
	Likelihood  int      `json:"likelihood"` // 1-5 scale
	Mitigation  string   `json:"mitigation,omitempty"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// Barrier represents execution barriers
type Barrier struct {
	Type        string   `json:"type"` // regulation, supply, distribution, trust, tech
	Description string   `json:"description"`
	Weight      float64  `json:"weight"` // 0.0-1.0
	EvidenceIDs []string `json:"evidence_ids"`
}

// GraveyardCase represents a failed similar startup
type GraveyardCase struct {
	CompanyName string   `json:"company_name"`
	Description string   `json:"description"`
	FailureCause string  `json:"failure_cause"`
	Lessons     string   `json:"lessons"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// MarketAnalysis represents market size and competition analysis
type MarketAnalysis struct {
	Competitors []Competitor `json:"competitors"`
	MarketStage string       `json:"market_stage"` // early, growing, mature, declining
	Positioning string       `json:"positioning"`
	EvidenceIDs []string     `json:"evidence_ids"`
}

// ProblemAnalysis represents problem validation analysis
type ProblemAnalysis struct {
	PainPoints  []string `json:"pain_points"`
	Validation  string   `json:"validation"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// BarrierAnalysis represents execution barrier analysis
type BarrierAnalysis struct {
	Barriers    []Barrier `json:"barriers"`
	EvidenceIDs []string  `json:"evidence_ids"`
}

// ExecutionAnalysis represents execution complexity analysis
type ExecutionAnalysis struct {
	CapitalRequirement string   `json:"capital_requirement"`
	TalentRarity      string   `json:"talent_rarity"`
	IntegrationCount  int      `json:"integration_count"`
	Complexity        float64  `json:"complexity"` // composite score
	EvidenceIDs       []string `json:"evidence_ids"`
}

// RiskAnalysis represents risk assessment
type RiskAnalysis struct {
	Risks       []Risk   `json:"risks"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// GraveyardAnalysis represents analysis of failed similar companies
type GraveyardAnalysis struct {
	Cases       []GraveyardCase `json:"cases"`
	EvidenceIDs []string        `json:"evidence_ids"`
}

// Viability represents the final verdict
type Viability struct {
	OverallScore    float64 `json:"overall_score"` // 0-100
	MarketScore     float64 `json:"market_score"`
	ProblemScore    float64 `json:"problem_score"`
	BarrierScore    float64 `json:"barrier_score"`
	ExecutionScore  float64 `json:"execution_score"`
	RiskScore       float64 `json:"risk_score"`
	GraveyardScore  float64 `json:"graveyard_score"`
	Recommendation  string  `json:"recommendation"`
	KeyInsights     []string `json:"key_insights"`
	EvidenceIDs     []string `json:"evidence_ids"`
}

// Analysis represents the complete analysis result
type Analysis struct {
	ID            string             `json:"id"`
	Idea          IdeaInput          `json:"idea"`
	Market        MarketAnalysis     `json:"market"`
	Problem       ProblemAnalysis    `json:"problem"`
	Barriers      BarrierAnalysis    `json:"barriers"`
	Execution     ExecutionAnalysis  `json:"execution"`
	Risks         RiskAnalysis       `json:"risks"`
	Graveyard     GraveyardAnalysis  `json:"graveyard"`
	Verdict       Viability          `json:"verdict"`
	Evidence      []Evidence         `json:"evidence"`
	CreatedAt     time.Time          `json:"created_at"`
	Partial       bool               `json:"partial,omitempty"` // if analysis was incomplete
	Meta          json.RawMessage    `json:"meta,omitempty"`    // analyzer raw outputs and validation
}

// ApproxLocation represents geographic location for search context
type ApproxLocation struct {
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
}

// SearchQuery represents a web search query
type SearchQuery struct {
	Query    string `json:"query"`
	Intent   string `json:"intent"` // competitors, funding, regulation, postmortems
	Priority int    `json:"priority"`
}

// CacheEntry represents a cached search result
type CacheEntry struct {
	Hash      string          `json:"hash" db:"hash"`
	Query     string          `json:"query" db:"query"`
	Result    json.RawMessage `json:"result" db:"result"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	TTL       int             `json:"ttl_seconds" db:"ttl_seconds"`
}

// AnalysisRequest represents an API request for analysis
type AnalysisRequest struct {
	Idea    IdeaInput       `json:"idea"`
	Options *AnalysisOptions `json:"options,omitempty"`
}

// AnalysisOptions represents optional parameters for analysis
type AnalysisOptions struct {
	MaxEvidence int            `json:"max_evidence,omitempty"`
	Location    *ApproxLocation `json:"location,omitempty"`
	Timeout     *time.Duration  `json:"timeout,omitempty"`
}

// GetLocation returns the location or nil if not set
func (ao *AnalysisOptions) GetLocation() *ApproxLocation {
	if ao == nil {
		return nil
	}
	return ao.Location
}

// AnalysisResponse represents the API response for analysis creation
type AnalysisResponse struct {
	AnalysisID string `json:"analysis_id"`
	Status     string `json:"status"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
