package middleware

import (
	"net/http"
	"time"

	"github.com/gustavo000/goLibGustavo/pkg/limiter"
)

func RateLimit(requests int, duration time.Duration) Middleware {
	rateLimiter := limiter.New(requests, duration)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			if !rateLimiter.Allow(clientIP) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
