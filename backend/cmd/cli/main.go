package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"realitycheck/internal/analyzers"
	"realitycheck/internal/app"
	"realitycheck/internal/cache"
	"realitycheck/internal/config"
	"realitycheck/internal/evidence"
	"realitycheck/internal/llm"
	"realitycheck/internal/report"
	"realitycheck/internal/schema"
	"realitycheck/internal/score"
	"realitycheck/internal/search"
	"realitycheck/internal/store"
	"realitycheck/pkg/types"
)

func main() {
	var (
		title      = flag.String("title", "", "Startup title (required)")
		oneLiner   = flag.String("one-liner", "", "One-liner description (required)")
		category   = flag.String("category", "", "Optional category")
		location   = flag.String("location", "", "Optional location (country or region)")
		output     = flag.String("out", "", "Output file path (default: stdout)")
		format     = flag.String("format", "markdown", "Output format: markdown, html, json")
		timeout    = flag.Duration("timeout", 60*time.Second, "Analysis timeout")
		maxEvidence = flag.Int("max-evidence", 20, "Maximum evidence items to collect")
		dbDSN      = flag.String("db", "", "Database DSN (uses config if not provided)")
		help       = flag.Bool("help", false, "Show help message")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "RealityCheck CLI - Startup Idea Analysis Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --title \"Loom\" --one-liner \"Agentic coding assistant\" --out report.md\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --title \"TaskAI\" --one-liner \"AI task automation\" --format html --out report.html\n", os.Args[0])
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Validate required arguments
	if *title == "" || *oneLiner == "" {
		fmt.Fprintf(os.Stderr, "Error: --title and --one-liner are required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Validate format
	if *format != "markdown" && *format != "html" && *format != "json" {
		fmt.Fprintf(os.Stderr, "Error: --format must be one of: markdown, html, json\n")
		os.Exit(1)
	}

	// Load configuration
	cfg := config.Load()
	
	// Override database DSN if provided
	if *dbDSN != "" {
		cfg.DatabaseDSN = *dbDSN
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Run analysis
	result, err := runAnalysis(cfg, *title, *oneLiner, *category, *location, *timeout, *maxEvidence)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Generate output
	var content string
	switch *format {
	case "markdown":
		builder := report.NewMarkdownBuilder()
		content = builder.Build(result)
	case "html":
		builder := report.NewHTMLBuilder()
		content = builder.Build(result)
	case "json":
		content = formatJSON(result)
	}

	// Write output
	if err := writeOutput(content, *output); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	fmt.Printf("Analysis completed successfully. Overall score: %.1f/100\n", result.Verdict.OverallScore)
	if *output != "" {
		fmt.Printf("Report saved to: %s\n", *output)
	}
}

func runAnalysis(cfg *config.Config, title, oneLiner, category, location string, timeout time.Duration, maxEvidence int) (types.Analysis, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout+30*time.Second) // Add buffer for setup
	defer cancel()

	// Initialize database
	db, err := schema.InitDatabase(ctx, cfg.DatabaseDSN)
	if err != nil {
		return types.Analysis{}, fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Run migrations
	if err := schema.Migrate(ctx, db); err != nil {
		return types.Analysis{}, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize components
	llmClient := llm.NewClient(cfg.OpenAIAPIKey, cfg.OpenAIRPS, cfg.OpenAIBurst)
	
	evidenceCache, err := cache.NewEvidenceCache(db, cfg.CacheLRUSize, cfg.CacheTTL)
	if err != nil {
		return types.Analysis{}, fmt.Errorf("failed to initialize evidence cache: %w", err)
	}

	planner := search.NewPlanner(cfg.MaxQueries)
	executor := search.NewExecutor(llmClient, evidenceCache, cfg.AnalysisTimeout)
	normalizer := evidence.NewNormalizer()
	calculator := score.NewCalculator(nil) // Use default weights
	coordinator := analyzers.NewCoordinator(llmClient, calculator)
	repository := store.NewRepository(db)

	orchestrator := app.NewOrchestrator(
		planner,
		executor,
		normalizer,
		coordinator,
		repository,
		maxEvidence,
		timeout,
	)

	// Create analysis request
	idea := types.IdeaInput{
		Title:    title,
		OneLiner: oneLiner,
		Category: category,
		Location: location,
	}

	var analysisLocation *types.ApproxLocation
	if location != "" {
		analysisLocation = &types.ApproxLocation{
			Country: location,
		}
	}

	request := types.AnalysisRequest{
		Idea: idea,
		Options: &types.AnalysisOptions{
			MaxEvidence: maxEvidence,
			Location:    analysisLocation,
			Timeout:     &timeout,
		},
	}

	// Run analysis
	fmt.Printf("Analyzing startup idea: %s\n", title)
	fmt.Printf("Description: %s\n", oneLiner)
	fmt.Printf("Timeout: %v\n", timeout)
	fmt.Printf("Max evidence: %d\n", maxEvidence)
	fmt.Println()

	analysisID, err := orchestrator.AnalyzeIdea(ctx, request)
	if err != nil {
		return types.Analysis{}, fmt.Errorf("analysis failed: %w", err)
	}

	// Retrieve the completed analysis
	result, err := orchestrator.GetAnalysis(ctx, analysisID)
	if err != nil {
		return types.Analysis{}, fmt.Errorf("failed to retrieve analysis: %w", err)
	}

	return result, nil
}

func formatJSON(analysis types.Analysis) string {
	// For CLI output, we'll create a simplified JSON representation
	simplified := map[string]interface{}{
		"id":         analysis.ID,
		"idea":       analysis.Idea,
		"verdict":    analysis.Verdict,
		"created_at": analysis.CreatedAt,
		"partial":    analysis.Partial,
		"scores": map[string]float64{
			"overall":   analysis.Verdict.OverallScore,
			"market":    analysis.Verdict.MarketScore,
			"problem":   analysis.Verdict.ProblemScore,
			"barriers":  analysis.Verdict.BarrierScore,
			"execution": analysis.Verdict.ExecutionScore,
			"risks":     analysis.Verdict.RiskScore,
			"graveyard": analysis.Verdict.GraveyardScore,
		},
		"evidence_count": len(analysis.Evidence),
	}

	// Use a simple JSON format for CLI output
	bytes, _ := json.MarshalIndent(simplified, "", "  ")
	return string(bytes)
}

func writeOutput(content, outputPath string) error {
	if outputPath == "" {
		// Write to stdout
		fmt.Print(content)
		return nil
	}

	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Write to file
	return os.WriteFile(outputPath, []byte(content), 0644)
}
