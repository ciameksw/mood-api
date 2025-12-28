package server

import (
	"net/http"

	"github.com/ciameksw/mood-api/mood/internal/mood/config"
	"github.com/ciameksw/mood-api/mood/internal/mood/postgres"
	"github.com/ciameksw/mood-api/pkg/logger"
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

	r.HandleFunc("POST /mood", s.handleAddMood)
	r.HandleFunc("GET /mood", s.handleGetMoods)
	r.HandleFunc("GET /mood/types", s.handleGetMoodTypes)
	r.HandleFunc("GET /mood/summary", s.handleGetMoodSummary)
	r.HandleFunc("PUT /mood", s.handleUpdateMood)
	r.HandleFunc("DELETE /mood/{id}", s.handleDeleteMood)

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
