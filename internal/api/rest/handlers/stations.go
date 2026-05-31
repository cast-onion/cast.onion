package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cast-onion/internal/service"
	"github.com/go-chi/chi/v5"
)

type StationHandler struct {
	svc *service.StationService
}

func NewStationHandler(svc *service.StationService) *StationHandler {
	return &StationHandler{svc: svc}
}

func (h *StationHandler) Directory(w http.ResponseWriter, r *http.Request) {
	stations, err := h.svc.Directory(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stations)
}

func (h *StationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	station, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(station)
}
