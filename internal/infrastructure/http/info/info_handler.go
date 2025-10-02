package info

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/MIREASHKI-BIG-BOB/backend_main/internal/infrastructure/ports/repository"
)

type Handler struct {
	infoRepo repository.InfoRepository
}

func New(infoRepo repository.InfoRepository) *Handler {
	return &Handler{
		infoRepo: infoRepo,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) GetDoctor(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "id parameter is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "invalid id parameter"})
		return
	}

	doctor, err := h.infoRepo.GetDoctorByID(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, doctor)
}

func (h *Handler) GetMedical(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "id parameter is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "invalid id parameter"})
		return
	}

	medical, err := h.infoRepo.GetMedicalByID(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, medical)
}

