package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MIREASHKI-BIG-BOB/backend_main/config"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/adapters/websocket/sensors"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/domain/services"
	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/database"
	healthHandler "github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/http/health"
)

type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *database.DB

	healthService *services.HealthService
	healthHandler *healthHandler.Handler
	sensorHandler *sensors.Handler

	router *chi.Mux
	server *http.Server
}

func New(cfg *config.Config, _ *database.DB) (*Server, error) {
	s := &Server{
		cfg:    cfg,
		logger: slog.Default(),
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("init server: %w", err)
	}

	return s, nil
}

func (s *Server) init() error {
	if err := s.initDB(); err != nil {
		return fmt.Errorf("init db: %v", err)
	}

	s.initServices()
	s.initHandlers()
	s.initRouter()
	s.initHTTPServer()

	return nil
}

func (s *Server) initDB() error {
	ctx := context.Background()
	db, err := database.Connect(ctx, &database.Config{
		Driver: s.cfg.DB.Driver,
		DSN:    s.cfg.DB.DSN,
	})
	if err != nil {
		return err
	}
	s.db = db

	if err := database.Migrate(ctx, &database.Config{
		Driver: s.cfg.DB.Driver,
		DSN:    s.cfg.DB.DSN,
	}); err != nil {
		return fmt.Errorf("migrate db: %v", err)
	}

	return nil
}

func (s *Server) initServices() {
	s.healthService = services.NewHealthService()
}

func (s *Server) initHandlers() {
	s.healthHandler = healthHandler.New(s.healthService)
	s.initSensorHandlers()
}

func (s *Server) initSensorHandlers() {
	allowedSensorsToToken := make(map[string]string, len(s.cfg.Sensors.Entities))
	for _, sensor := range s.cfg.Sensors.Entities {
		allowedSensorsToToken[sensor.UUID] = sensor.Token
	}

	sensorCfg := &sensors.Config{
		AllowedSensorsToToken: allowedSensorsToToken,
		HandshakeTimeout:      s.cfg.Sensors.HandshakeTimeout,
	}
	s.sensorHandler = sensors.NewHandler(sensorCfg, s.logger)
}

func (s *Server) initHTTPServer() {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Addr, s.cfg.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
	}
}

func (s *Server) Run() error {
	s.logger.Info("Server running", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server listen: %w", err)
	}

	return nil
}
