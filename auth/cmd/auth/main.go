package main

import (
	"github.com/ciameksw/mood-api/auth/internal/auth/config"
	"github.com/ciameksw/mood-api/auth/internal/auth/logger"
	"github.com/ciameksw/mood-api/auth/internal/auth/postgres"
	"github.com/ciameksw/mood-api/auth/internal/auth/server"
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
	defer db.Disconnect()

	s := server.NewServer(lgr, cfg, db)
	s.Start()
}
