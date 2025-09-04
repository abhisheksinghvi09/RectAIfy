package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"realitycheck/internal/analyzers"
	"realitycheck/internal/app"
	"realitycheck/internal/cache"
	"realitycheck/internal/config"
	"realitycheck/internal/evidence"
	"realitycheck/internal/llm"
	"realitycheck/internal/schema"
	"realitycheck/internal/score"
	"realitycheck/internal/search"
	"realitycheck/internal/store"
	"realitycheck/pkg/httpx"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Initialize database
	ctx := context.Background()
	db, err := schema.InitDatabase(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := schema.Migrate(ctx, db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize components
	llmClient := llm.NewClient(cfg.OpenAIAPIKey, cfg.OpenAIRPS, cfg.OpenAIBurst)

	evidenceCache, err := cache.NewEvidenceCache(db, cfg.CacheLRUSize, cfg.CacheTTL)
	if err != nil {
		log.Fatalf("Failed to initialize evidence cache: %v", err)
	}

	// Start cache cleanup worker
	go evidenceCache.StartCleanupWorker(ctx, time.Hour)

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
		cfg.MaxEvidencePerQuery,
		cfg.AnalysisTimeout,
	)

	// Initialize HTTP handlers
	handlers := httpx.NewAPIHandlers(orchestrator)

	// Setup HTTP server
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/v1/analyze", handlers.HandleAnalyze)
	mux.HandleFunc("/v1/analyses/", handlers.HandleGetAnalysis)
	mux.HandleFunc("/v1/analyses", handlers.HandleListAnalyses)
	mux.HandleFunc("/v1/stats", handlers.HandleStats)
	mux.HandleFunc("/health", handlers.HandleHealthCheck)

	// Apply middleware
	var handler http.Handler = mux
	handler = httpx.AuthMiddleware(cfg.BearerToken)(handler)
	handler = httpx.LoggingMiddleware(handler)
	handler = httpx.CORSMiddleware(handler)

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: handler,

		// Timeouts
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second, // Long timeout for analysis requests
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting RealityCheck API server on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
