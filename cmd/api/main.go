package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kimbef/19-dec-2025/internal/config"
	"github.com/kimbef/19-dec-2025/internal/database"
	"github.com/kimbef/19-dec-2025/internal/handlers"
	"github.com/kimbef/19-dec-2025/internal/repository"
	"github.com/kimbef/19-dec-2025/internal/server"
	"github.com/kimbef/19-dec-2025/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	log, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	log.Info("Starting application...")

	// Initialize database
	db, err := database.New(database.Config{
		DSN:             cfg.Database.GetDSN(),
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	log.Info("Connected to database successfully")

	// Run migrations
	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	log.Info("Database migrations completed")

	// Initialize repositories
	itemRepo := repository.NewItemRepository(db.DB)

	// Initialize handlers
	h := handlers.New(itemRepo, log)

	// Initialize server
	srv := server.New(server.Config{
		Port:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Handler:      h,
		Logger:       log,
	})

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Application stopped")
}
