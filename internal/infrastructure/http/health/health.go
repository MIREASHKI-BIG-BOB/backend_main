package health

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/services"
)

type Handler struct {
	healthService *services.HealthService
}

func New(
	healthService *services.HealthService,
) *Handler {
	return &Handler{
		healthService: healthService,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	healthInfo := h.healthService.GetHealthStatus()

	render.JSON(w, r, healthInfo)
}
