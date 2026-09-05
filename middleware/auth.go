package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwtauth "github.com/gustavo000/goLibGustavo/pkg/auth"
)

type contextKey string

const (
	UserKey   contextKey = "user"
	TokenKey  contextKey = "token"
	ClaimsKey contextKey = "claims"
)

func Auth2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
		}
	})
}

func Auth(next http.Handler) http.Handler {
	return AuthWithSecret(os.Getenv("JWT_SECRET"))(next)
}

func AuthWithSecret(secret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for public endpoints
			if isPublicEndpoint(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			claims, err := jwtauth.ParseToken(token, secret)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, claims.UserID)
			ctx = context.WithValue(ctx, TokenKey, token)
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicEndpoint(path string) bool {
	return path == "/login" || path == "/health"
}
