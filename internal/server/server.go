package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kimbef/19-dec-2025/internal/handlers"
	"github.com/kimbef/19-dec-2025/internal/middleware"
	"github.com/kimbef/19-dec-2025/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	log        *logger.Logger
}

// Config holds server configuration
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Handler      *handlers.Handler
	Logger       *logger.Logger
}

// New creates a new server instance
func New(cfg Config) *Server {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", cfg.Handler.Health)

	// API routes
	mux.HandleFunc("POST /api/items", cfg.Handler.CreateItem)
	mux.HandleFunc("GET /api/items/{id}", cfg.Handler.GetItem)
	mux.HandleFunc("PUT /api/items/{id}", cfg.Handler.UpdateItem)
	mux.HandleFunc("DELETE /api/items/{id}", cfg.Handler.DeleteItem)
	mux.HandleFunc("GET /api/items", cfg.Handler.ListItems)

	// Apply middleware
	handler := middleware.Chain(
		middleware.Logger(cfg.Logger),
		middleware.CORS(),
		middleware.Recovery(cfg.Logger),
	)(mux)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
		log:        cfg.Logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.log.Infof("Starting server on %s", s.httpServer.Addr)
	
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.log.Info("Server stopped gracefully")
	return nil
}
