package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/MIREASHKI-BIG-BOB/backend_main/config"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/domain/services"
	healthHandler "github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/http/health"
)

type Server struct {
	cfg *config.Config

	// services
	healthService *services.HealthService

	// handlers
	healthHandler *healthHandler.Handler

	// infrastructure
	router *chi.Mux
	server *http.Server
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		cfg: cfg,
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("init server: %w", err)
	}

	return s, nil
}

func (s *Server) init() error {
	s.initServices()
	s.initHandlers()
	s.initRouter()
	s.initHTTPServer()

	return nil
}

func (s *Server) initServices() {
	s.healthService = services.NewHealthService()
}

func (s *Server) initHandlers() {
	s.healthHandler = healthHandler.New(s.healthService)
}

func (s *Server) initHTTPServer() {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Addr, s.cfg.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server listen: %w", err)
	}

	return nil
}
