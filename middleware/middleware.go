package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gustavo000/goLibGustavo/models"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return // Short-circuit: Stop the chain here
		}

		// Validate token and get user info
		userID := "12345"

		// Pass userID to the next handler via Context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
