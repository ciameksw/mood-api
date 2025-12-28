package server

import (
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
}

func NewServer(log *logger.Logger, cfg *config.Config, pg *postgres.PostgresDB) *Server {
	return &Server{
		Logger:       log,
		Config:       cfg,
		DBOperations: &repository.DBOperations{Postgres: pg},
		Validator:    validator.New(),
	}
}

func (s *Server) Start() {
	r := http.NewServeMux()

	r.HandleFunc("POST /advice/select", s.handleSelectAdvice)

	r.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := s.Config.ServerHost + ":" + s.Config.ServerPort
	s.Logger.Info.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		s.Logger.Error.Fatalf("Server failed to start: %v", err)
	}
}
