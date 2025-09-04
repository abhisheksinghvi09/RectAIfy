package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"rectaify/pkg/types"
)

// Client wraps OpenAI API with rate limiting and web search
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	limiter    *rate.Limiter
}

// NewClient creates a new OpenAI client with rate limiting
func NewClient(apiKey string, rps int, burst int) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

// SearchRequest represents a web search request
type SearchRequest struct {
	Model    string              `json:"model"`
	Messages []ChatMessage       `json:"messages"`
	Tools    []Tool              `json:"tools"`
	ToolChoice string            `json:"tool_choice"`
	Temperature float64          `json:"temperature"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool represents a tool definition
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents a function tool
type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters,omitempty"`
}

// SearchResponse represents the OpenAI response
type SearchResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a response choice
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	ToolCalls    []ToolCall  `json:"tool_calls,omitempty"`
	FinishReason string      `json:"finish_reason"`
}

// ToolCall represents a tool call
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// WebSearchResult represents a web search result
type WebSearchResult struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// Search performs web search using OpenAI's web_search_preview
func (c *Client) Search(ctx context.Context, queries []string, location *types.ApproxLocation) ([]types.Evidence, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	var evidence []types.Evidence
	
	for _, query := range queries {
		results, err := c.performWebSearch(ctx, query, location)
		if err != nil {
			// Log error but continue with other queries
			continue
		}
		
		for _, result := range results {
			ev := types.Evidence{
				ID:          generateEvidenceID(result.URL, result.Title, result.PublishedAt),
				URL:         result.URL,
				Title:       result.Title,
				Snippet:     result.Content,
				PublishedAt: result.PublishedAt,
				RetrievedAt: time.Now(),
				SourceType:  inferSourceType(result.URL),
			}
			evidence = append(evidence, ev)
		}
	}

	return evidence, nil
}

// ConstrainedJSON performs a constrained JSON generation request
func (c *Client) ConstrainedJSON(ctx context.Context, systemPrompt string, userPrompt interface{}, schema []byte) (json.RawMessage, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	// Convert user prompt to string if needed
	var userString string
	switch v := userPrompt.(type) {
	case string:
		userString = v
	default:
		userBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user prompt: %w", err)
		}
		userString = string(userBytes)
	}

	// Parse schema for response format
	var schemaObj map[string]interface{}
	if err := json.Unmarshal(schema, &schemaObj); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	request := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userString},
		},
		"temperature": 0.2,
		"response_format": map[string]interface{}{
			"type":        "json_schema",
			"json_schema": map[string]interface{}{
				"name":   "analysis_response",
				"strict": true,
				"schema": schemaObj,
			},
		},
	}

	response, err := c.makeRequest(ctx, "/chat/completions", request)
	if err != nil {
		return nil, err
	}

	var chatResponse SearchResponse
	if err := json.Unmarshal(response, &chatResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResponse.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	return json.RawMessage(chatResponse.Choices[0].Message.Content), nil
}

// performWebSearch executes a web search query
func (c *Client) performWebSearch(ctx context.Context, query string, location *types.ApproxLocation) ([]WebSearchResult, error) {
	locationStr := ""
	if location != nil && location.Country != "" {
		locationStr = fmt.Sprintf(" in %s", location.Country)
		if location.Region != "" {
			locationStr = fmt.Sprintf(" in %s, %s", location.Region, location.Country)
		}
	}

	searchQuery := query + locationStr

	request := SearchRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("Search for information about: %s", searchQuery),
			},
		},
		Tools: []Tool{
			{
				Type: "web_search",
				Function: ToolFunction{
					Name:        "web_search",
					Description: "Search the web for current information",
				},
			},
		},
		ToolChoice:  "required",
		Temperature: 0.2,
	}

	response, err := c.makeRequest(ctx, "/chat/completions", request)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchResponse
	if err := json.Unmarshal(response, &searchResponse); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	// Extract web search results from tool calls
	var results []WebSearchResult
	for _, choice := range searchResponse.Choices {
		for _, toolCall := range choice.ToolCalls {
			if toolCall.Function.Name == "web_search" {
				var searchResults []WebSearchResult
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &searchResults); err != nil {
					continue // Skip malformed results
				}
				results = append(results, searchResults...)
			}
		}
	}

	return results, nil
}

// makeRequest performs an HTTP request to the OpenAI API
func (c *Client) makeRequest(ctx context.Context, endpoint string, payload interface{}) ([]byte, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
