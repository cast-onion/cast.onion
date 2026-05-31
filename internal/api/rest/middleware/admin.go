package middleware

import (
	"context"
	"net/http"

	"github.com/cast-onion/internal/db/queries"
	"github.com/cast-onion/pkg/keygen"
)

type adminKey string

const AdminKey adminKey = "admin_id"

func RequireAdmin(q *queries.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := r.Header.Get("X-Admin-Token")
			if raw == "" {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			hash := keygen.HashRaw(raw)
			adminID, err := q.GetAdminIDByTokenHash(r.Context(), hash)
			if err != nil {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), AdminKey, adminID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminFromContext(ctx context.Context) string {
	v, _ := ctx.Value(AdminKey).(string)
	return v
}
