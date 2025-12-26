package server

import (
	"net/http"

	"github.com/ciameksw/mood-api/auth/internal/auth/config"
	"github.com/ciameksw/mood-api/auth/internal/auth/logger"
	"github.com/ciameksw/mood-api/auth/internal/auth/postgres"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Logger    *logger.Logger
	Config    *config.Config
	Postgres  *postgres.PostgresDB
	Validator *validator.Validate
}

func NewServer(log *logger.Logger, cfg *config.Config, pg *postgres.PostgresDB) *Server {
	return &Server{
		Logger:    log,
		Config:    cfg,
		Postgres:  pg,
		Validator: validator.New(),
	}
}

func (s *Server) Start() {
	r := http.NewServeMux()

	r.HandleFunc("POST /auth/login", s.handleLogin)
	r.HandleFunc("POST /auth/register", s.handleRegister)
	r.HandleFunc("GET /auth/authorize", s.handleAuthorize)

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
