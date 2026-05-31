package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cast-onion/internal/api/rest/middleware"
	"github.com/cast-onion/internal/service"
	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	apps, err := h.svc.ListPending(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(apps)
}

func (h *AdminHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	adminID := middleware.AdminFromContext(r.Context())
	result, err := h.svc.ApproveApplication(r.Context(), adminID, id, body.Reason)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyReviewed) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (h *AdminHandler) Deny(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	adminID := middleware.AdminFromContext(r.Context())
	if err := h.svc.DenyApplication(r.Context(), adminID, id, body.Reason); err != nil {
		if errors.Is(err, service.ErrAlreadyReviewed) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	adminID := middleware.AdminFromContext(r.Context())
	if err := h.svc.SuspendStation(r.Context(), adminID, id, body.Reason); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	adminID := middleware.AdminFromContext(r.Context())
	if err := h.svc.RevokeStation(r.Context(), adminID, id, body.Reason); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) Unsuspend(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	adminID := middleware.AdminFromContext(r.Context())
	if err := h.svc.UnsuspendStation(r.Context(), adminID, id, body.Reason); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
