package server

import (
	"context"
	"net/http"

	"github.com/ciameksw/mood-api/gateway/internal/gateway/config"
	"github.com/ciameksw/mood-api/gateway/internal/gateway/services/auth"
	"github.com/ciameksw/mood-api/gateway/internal/gateway/services/mood"
	"github.com/ciameksw/mood-api/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Logger      *logger.Logger
	Config      *config.Config
	AuthService *auth.AuthService
	MoodService *mood.MoodService
	Validator   *validator.Validate
	httpServer  *http.Server
}

func NewServer(log *logger.Logger, cfg *config.Config) *Server {
	return &Server{
		Logger:      log,
		Config:      cfg,
		AuthService: auth.NewAuthService(cfg),
		MoodService: mood.NewMoodService(cfg),
		Validator:   validator.New(),
	}
}

func (s *Server) Start() error {
	r := http.NewServeMux()

	s.setupAuthRouter(r)
	s.setupMoodRouter(r)
	s.setupAdviceRouter(r)
	s.setupQuoteRouter(r)

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
