package handlers

import (
	"net/http"

	"github.com/cast-onion/internal/service"
)

type StreamHandler struct {
	svc *service.StreamService
}

func NewStreamHandler(svc *service.StreamService) *StreamHandler {
	return &StreamHandler{svc: svc}
}

func (h *StreamHandler) Auth(w http.ResponseWriter, r *http.Request) {
	rawKey := r.FormValue("pass")
	if err := h.svc.ValidateStationKey(r.Context(), rawKey); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.Header().Set("icecast-auth-user", "1")
	w.WriteHeader(http.StatusOK)
}
