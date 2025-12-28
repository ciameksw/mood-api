package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ciameksw/mood-api/auth/internal/auth/config"
	"github.com/ciameksw/mood-api/auth/internal/auth/server"
	"github.com/ciameksw/mood-api/pkg/logger"
	"github.com/ciameksw/mood-api/pkg/postgres"
)

func main() {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()

	// Connect to Postgres
	db, err := postgres.Connect(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDatabase, cfg.PostgresSSLMode)
	if err != nil {
		lgr.Error.Fatalf("Failed to connect to Postgres: %v", err)
	}

	s := server.NewServer(lgr, cfg, db)

	// Start server in a goroutine
	go func() {
		lgr.Info.Printf("Starting server on %s:%s", cfg.ServerHost, cfg.ServerPort)
		if err := s.Start(); err != nil {
			lgr.Error.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	lgr.Info.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		lgr.Error.Printf("Server forced to shutdown: %v", err)
	}
	db.Disconnect(ctx)

	lgr.Info.Println("Server exited gracefully")
}
