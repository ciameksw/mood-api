package server

import (
	"context"
	"net/http"

	"github.com/ciameksw/mood-api/advice/internal/advice/config"
	"github.com/ciameksw/mood-api/advice/internal/advice/repository"
	"github.com/ciameksw/mood-api/pkg/logger"
	"github.com/ciameksw/mood-api/pkg/postgres"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Logger       *logger.Logger
	Config       *config.Config
	DBOperations *repository.DBOperations
	Validator    *validator.Validate
	httpServer   *http.Server
}

func NewServer(log *logger.Logger, cfg *config.Config, pg *postgres.PostgresDB) *Server {
	return &Server{
		Logger:       log,
		Config:       cfg,
		DBOperations: &repository.DBOperations{Postgres: pg},
		Validator:    validator.New(),
	}
}

func (s *Server) Start() error {
	r := http.NewServeMux()

	r.HandleFunc("POST /advice/select", s.handleSelectAdvice)
	r.HandleFunc("POST /advice/period/save", s.handleSaveAdvice)
	r.HandleFunc("GET /advice/{id}", s.handleGetByID)

	r.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := s.Config.ServerHost + ":" + s.Config.ServerPort
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	s.Logger.Info.Printf("Starting server on %s", addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
