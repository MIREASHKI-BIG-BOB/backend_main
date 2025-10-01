package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func (s *Server) initRouter() {
	s.router = chi.NewRouter()

	// Middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// Routing
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/health", s.healthHandler.HealthCheck)

		r.Route("/sensors", func(r chi.Router) {
			r.Get("/start", s.sensorsHandler.StartSensor)
			r.Get("/stop", s.sensorsHandler.StopAllSensors)
		})
	})

	s.router.Route("/ws", func(r chi.Router) {
		r.HandleFunc("/sensor", s.sensorHandler.HandleWebSocket)
		r.HandleFunc("/", s.frontendHandler.HandleWebSocket)
	})
}
