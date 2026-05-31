package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cast-onion/internal/api/rest/middleware"
	"github.com/cast-onion/internal/email"
	"github.com/cast-onion/internal/service"
	"github.com/go-chi/chi/v5"
)

type ApplicationHandler struct {
	svc    *service.ApplicationService
	mailer *email.Dispatcher
}

func NewApplicationHandler(svc *service.ApplicationService, mailer *email.Dispatcher) *ApplicationHandler {
	return &ApplicationHandler{svc: svc, mailer: mailer}
}

func (h *ApplicationHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ContactEmail string `json:"contact_email"`
		StationName  string `json:"station_name"`
		Description  string `json:"description"`
		Genre        string `json:"genre"`
		Notes        string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	sessionID := middleware.SessionFromContext(r.Context())
	app, err := h.svc.Submit(r.Context(), sessionID, body.ContactEmail, body.StationName, body.Description, body.Genre, body.Notes)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.mailer.ApplicationReceived(
		body.ContactEmail,
		body.ContactEmail,
		body.StationName,
		body.Genre,
		app.ID,
	)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func (h *ApplicationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	app, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(app)
}
