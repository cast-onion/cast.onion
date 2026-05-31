package middleware

import (
	"context"
	"net/http"

	"github.com/cast-onion/internal/service"
)

type contextKey string

const SessionKey contextKey = "session_id"

func RequireSession(sessionSvc *service.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Session-ID")
			if id == "" {
				http.Error(w, "missing session", http.StatusUnauthorized)
				return
			}
			ok, err := sessionSvc.Validate(r.Context(), id)
			if err != nil || !ok {
				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), SessionKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SessionFromContext(ctx context.Context) string {
	v, _ := ctx.Value(SessionKey).(string)
	return v
}
