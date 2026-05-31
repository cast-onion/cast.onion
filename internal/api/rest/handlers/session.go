package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cast-onion/internal/api/rest/middleware"
)

func SessionInfo(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.SessionFromContext(r.Context())
	json.NewEncoder(w).Encode(map[string]string{
		"session_id": sessionID,
	})
}
