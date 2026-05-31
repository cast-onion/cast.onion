package middleware

import (
	"net/http"
	"time"

	"github.com/cast-onion/internal/cache"
)

func RateLimit(redis *cache.Redis, limit int64, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr
			count, err := redis.IncrRateLimit(r.Context(), key, window)
			if err != nil || count > limit {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
