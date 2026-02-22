package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gustavo000/goLibGustavo/pkg/response"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &response.ResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Printf(
			"[%s] %s %s %d %s | RequestID: %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			wrapped.StatusCode,
			time.Since(start),
			r.Context().Value(RequestIDKey),
		)
	})
}
