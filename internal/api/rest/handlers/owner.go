package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cast-onion/internal/service"
)

type OwnerHandler struct {
	authSvc    *service.AuthService
	stationSvc *service.StationService
}

func NewOwnerHandler(authSvc *service.AuthService, stationSvc *service.StationService) *OwnerHandler {
	return &OwnerHandler{authSvc: authSvc, stationSvc: stationSvc}
}

func (h *OwnerHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	raw := r.Header.Get("X-Access-Token")
	if raw == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	stationID, err := h.authSvc.ValidateAccessToken(r.Context(), raw)
	if err != nil {
		log.Printf("owner dashboard: token validation failed: %v", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	log.Printf("owner dashboard: token valid, station_id=%s", stationID)
	station, err := h.stationSvc.GetByID(r.Context(), stationID)
	if err != nil {
		log.Printf("owner dashboard: station not found for id=%s err=%v", stationID, err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)
}

func (h *OwnerHandler) UpdateStation(w http.ResponseWriter, r *http.Request) {
	raw := r.Header.Get("X-Access-Token")
	if raw == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	stationID, err := h.authSvc.ValidateAccessToken(r.Context(), raw)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	var body struct {
		Description string `json:"description"`
		Genre       string `json:"genre"`
		WebsiteURL  string `json:"website_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := h.stationSvc.Update(r.Context(), stationID, body.Description, body.Genre, body.WebsiteURL); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
