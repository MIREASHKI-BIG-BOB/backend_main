package sensors

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/usecases/sensors"
)

type Handler struct {
	sensorsUseCase sensors.SensorsUseCase
}

func New(sensorsUseCase sensors.SensorsUseCase) *Handler {
	return &Handler{
		sensorsUseCase: sensorsUseCase,
	}
}

type StartSensorResponse struct {
	Message string                `json:"message"`
	Sensor  *sensors.SensorStatus `json:"sensor"`
}

type StopSensorsResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) StartSensor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sensor, err := h.sensorsUseCase.ConnectSensor(ctx)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, StartSensorResponse{
		Message: "Sensor started successfully",
		Sensor:  sensor,
	})
}
func (h *Handler) StopAllSensors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.sensorsUseCase.DisconnectSensors(ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, StopSensorsResponse{
		Message: "All sensors stopped successfully",
	})
}
