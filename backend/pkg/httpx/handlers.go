package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"realitycheck/internal/app"
	"realitycheck/internal/report"
	"realitycheck/pkg/types"
)

// APIHandlers contains all HTTP handlers for the API
type APIHandlers struct {
	orchestrator    *app.Orchestrator
	markdownBuilder *report.MarkdownBuilder
	htmlBuilder     *report.HTMLBuilder
}

// NewAPIHandlers creates new API handlers
func NewAPIHandlers(orchestrator *app.Orchestrator) *APIHandlers {
	return &APIHandlers{
		orchestrator:    orchestrator,
		markdownBuilder: report.NewMarkdownBuilder(),
		htmlBuilder:     report.NewHTMLBuilder(),
	}
}

// HandleAnalyze handles POST /v1/analyze
func (h *APIHandlers) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request types.AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.Idea.Title == "" || request.Idea.OneLiner == "" {
		h.writeErrorResponse(w, "Title and OneLiner are required", http.StatusBadRequest)
		return
	}

	// Start analysis
	analysisID, err := h.orchestrator.AnalyzeIdea(r.Context(), request)
	if err != nil {
		h.writeErrorResponse(w, fmt.Sprintf("Analysis failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := types.AnalysisResponse{
		AnalysisID: analysisID,
		Status:     "completed",
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// HandleGetAnalysis handles GET /v1/analyses/{id}
func (h *APIHandlers) HandleGetAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract analysis ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/v1/analyses/")
	analysisID := strings.Split(path, ".")[0] // Remove file extension if present

	if analysisID == "" {
		h.writeErrorResponse(w, "Analysis ID is required", http.StatusBadRequest)
		return
	}

	analysis, err := h.orchestrator.GetAnalysis(r.Context(), analysisID)
	if err != nil {
		if err.Error() == "analysis not found" {
			h.writeErrorResponse(w, "Analysis not found", http.StatusNotFound)
			return
		}
		h.writeErrorResponse(w, fmt.Sprintf("Failed to get analysis: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if a specific format is requested
	if strings.HasSuffix(r.URL.Path, ".md") {
		h.handleMarkdownResponse(w, analysis)
		return
	}

	if strings.HasSuffix(r.URL.Path, ".html") {
		h.handleHTMLResponse(w, analysis)
		return
	}

	// Default to JSON
	h.writeJSONResponse(w, analysis, http.StatusOK)
}

// HandleListAnalyses handles GET /v1/analyses
func (h *APIHandlers) HandleListAnalyses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	searchQuery := r.URL.Query().Get("q")

	limit := 10 // default
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	var analyses []types.Analysis
	var err error

	if searchQuery != "" {
		analyses, err = h.orchestrator.SearchAnalyses(r.Context(), searchQuery, limit, offset)
	} else {
		analyses, err = h.orchestrator.ListAnalyses(r.Context(), limit, offset)
	}

	if err != nil {
		h.writeErrorResponse(w, fmt.Sprintf("Failed to list analyses: %v", err), http.StatusInternalServerError)
		return
	}

	// Create response with pagination info
	totalCount, _ := h.orchestrator.GetAnalysisCount(r.Context())
	
	response := map[string]interface{}{
		"analyses": analyses,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"total":  totalCount,
		},
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// HandleDeleteAnalysis handles DELETE /v1/analyses/{id}
func (h *APIHandlers) HandleDeleteAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract analysis ID from URL path
	analysisID := strings.TrimPrefix(r.URL.Path, "/v1/analyses/")

	if analysisID == "" {
		h.writeErrorResponse(w, "Analysis ID is required", http.StatusBadRequest)
		return
	}

	err := h.orchestrator.DeleteAnalysis(r.Context(), analysisID)
	if err != nil {
		if err.Error() == "analysis not found" {
			h.writeErrorResponse(w, "Analysis not found", http.StatusNotFound)
			return
		}
		h.writeErrorResponse(w, fmt.Sprintf("Failed to delete analysis: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleHealthCheck handles GET /health
func (h *APIHandlers) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.orchestrator.HealthCheck(r.Context())
	if err != nil {
		h.writeErrorResponse(w, fmt.Sprintf("Health check failed: %v", err), http.StatusServiceUnavailable)
		return
	}

	response := map[string]string{
		"status": "healthy",
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// HandleStats handles GET /v1/stats
func (h *APIHandlers) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.orchestrator.GetStats(r.Context())
	if err != nil {
		h.writeErrorResponse(w, fmt.Sprintf("Failed to get stats: %v", err), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, stats, http.StatusOK)
}

// handleMarkdownResponse sends analysis as markdown
func (h *APIHandlers) handleMarkdownResponse(w http.ResponseWriter, analysis types.Analysis) {
	markdown := h.markdownBuilder.Build(analysis)
	
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.md\"", analysis.ID))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(markdown))
}

// handleHTMLResponse sends analysis as HTML
func (h *APIHandlers) handleHTMLResponse(w http.ResponseWriter, analysis types.Analysis) {
	html := h.htmlBuilder.Build(analysis)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// writeJSONResponse writes a JSON response
func (h *APIHandlers) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If we can't encode the response, there's not much we can do
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// writeErrorResponse writes an error response
func (h *APIHandlers) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := types.ErrorResponse{
		Error: message,
	}
	h.writeJSONResponse(w, errorResponse, statusCode)
}
