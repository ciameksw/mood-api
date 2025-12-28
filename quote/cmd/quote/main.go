package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ciameksw/mood-api/pkg/logger"
	"github.com/ciameksw/mood-api/quote/internal/quote/config"
	"github.com/ciameksw/mood-api/quote/internal/quote/server"
)

func main() {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()

	s := server.NewServer(lgr, cfg)

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

	lgr.Info.Println("Server exited gracefully")
}
