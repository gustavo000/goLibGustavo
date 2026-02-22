package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserKey contextKey = "user"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		if r.URL.Path == "/login" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		/*
			if !strings.HasPrefix(token, "token-") {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}*/

		ctx := context.WithValue(r.Context(), UserKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
