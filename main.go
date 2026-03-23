package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UcGeorge/pastel-fde-assessment/handlers"
	"github.com/UcGeorge/pastel-fde-assessment/internal/config"
	"github.com/UcGeorge/pastel-fde-assessment/internal/di"
	"github.com/UcGeorge/pastel-fde-assessment/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
)

func main() {
	// Setup structured, human-readable logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Initialize the DI container
	injector := di.NewContainer()
	log.Info().Msg("Dependency injection container initialized")

	// Resolve configuration
	cfg := do.MustInvoke[*config.Config](injector)

	// Parse all templates at startup
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	log.Info().Int("templates", len(tmpl.Templates())).Msg("Templates parsed successfully")

	// Resolve services from DI container
	txnService := do.MustInvoke[*services.TransactionService](injector)
	screenService := do.MustInvoke[*services.ScreeningService](injector)
	adverseService := do.MustInvoke[*services.AdverseMediaService](injector)

	// Create handlers with injected dependencies
	homeHandler := handlers.NewHomeHandler(tmpl)
	txnHandler := handlers.NewTransactionHandler(tmpl, txnService)
	screenHandler := handlers.NewScreeningHandler(tmpl, screenService)
	adverseHandler := handlers.NewAdverseMediaHandler(tmpl, adverseService)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Pages
	mux.HandleFunc("/", homeHandler.ServeHTTP)
	mux.HandleFunc("/transaction", txnHandler.Page)
	mux.HandleFunc("/screening", screenHandler.Page)
	mux.HandleFunc("/adverse-media", adverseHandler.Page)

	// HTMX API endpoints
	mux.HandleFunc("/transaction/submit", txnHandler.Submit)
	mux.HandleFunc("/screening/pep", screenHandler.CheckPEP)
	mux.HandleFunc("/screening/sanction", screenHandler.CheckSanction)
	mux.HandleFunc("/adverse-media/check", adverseHandler.Check)

	// Create the server
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("addr", addr).Msg("Starting Sigma Compliance Dashboard")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exited gracefully")
}
