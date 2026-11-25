package net_http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gustavo000/goLibGustavo/constants"
	"github.com/gustavo000/goLibGustavo/middleware"
	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
)

// Config holds server configuration
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns a default configuration   ++
func DefaultConfig() *Config {
	return &Config{
		Port:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// App represents the application
type App struct {
	config *Config
	router *http.ServeMux
}

// Error handler middleware
func (app *App) errorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.handleError(w, r, http.StatusInternalServerError, "Internal server error", err)
			}
		}()

		// Create a custom response writer to capture status code
		crw := &customResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next(crw, r)

		// Log error responses (status code >= 400)
		if crw.statusCode >= 400 {
			log.Printf("Error response: %s %s - Status: %d", r.Method, r.URL.Path, crw.statusCode)
		}
	}
}

// rateLimitEnvNames returns the env vars to try, most specific first:
// RATE_LIMIT_<SERVICE>_<KIND> (per microservice), then <KIND>_RATE_LIMIT.
func rateLimitEnvNames(kind string) []string {
	names := []string{}
	if props := properties.GetProperty(); props != nil {
		if service := sanitizeEnvSegment(props.GetName()); service != "" {
			names = append(names, fmt.Sprintf("RATE_LIMIT_%s_%s", service, kind))
		}
	}
	return append(names, kind+"_RATE_LIMIT")
}

// sanitizeEnvSegment makes a service name safe for use inside an env var name.
func sanitizeEnvSegment(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

// rateLimit resolves the per-minute cap. Precedence:
// service-scoped env var (RATE_LIMIT_<SERVICE>_<KIND>), generic env var
// (<KIND>_RATE_LIMIT), the value from environment/*.properties.json, and
// finally the package constant.
func rateLimit(kind string, configured, fallback int) int {
	for _, name := range rateLimitEnvNames(kind) {
		if raw := os.Getenv(name); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil && n > 0 {
				return n
			}
		}
	}
	if configured > 0 {
		return configured
	}
	return fallback
}

// configuredRateLimits returns the public/private caps from the properties
// file, or zero when properties were not initialized.
func configuredRateLimits() (public, private int) {
	if props := properties.GetProperty(); props != nil {
		return props.GetPublicRateLimit(), props.GetPrivateRateLimit()
	}
	return 0, 0
}

// Start the server
func StartHttp(routes rest.Routes, port string) {
	// Create router
	mux := http.NewServeMux()
	// Public routes stay open, so the same visitor gets a tighter cap there.
	publicConfigured, privateConfigured := configuredRateLimits()
	publicLimit := middleware.RateLimit(rateLimit("PUBLIC", publicConfigured, constants.PUBLIC_RATE_LIMIT), time.Minute)
	methodsByPattern := map[string][]string{}
	for _, rt := range routes {
		route := rt
		routeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if route.Method != "" && r.Method != route.Method {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			resp := route.Controller.Service(r, &rest.Response{})
			if resp == nil || resp.GetHttp() == nil {
				http.Error(w, "empty response", http.StatusInternalServerError)
				return
			}
			for k, vals := range resp.GetHeader() {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", "application/json")
			}
			w.WriteHeader(resp.GetHttp().StatusCode)
			if body := resp.CopyBody(); body != nil {
				_, _ = io.Copy(w, body)
			}
		})

		var handler http.Handler = routeHandler
		if route.Public {
			handler = publicLimit(handler)
		} else {
			handler = middleware.Auth(handler)
		}
		mux.Handle(fmt.Sprintf("%s %s", route.Method, route.Pattern), handler)
		methodsByPattern[route.Pattern] = append(methodsByPattern[route.Pattern], route.Method)
	}
	for pattern, methods := range methodsByPattern {
		allow := middleware.AllowMethods(methods)
		mux.HandleFunc("OPTIONS "+pattern, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", allow)
			w.Header().Set("Allow", allow)
			w.WriteHeader(http.StatusNoContent)
		})
	}

	// Apply middleware chain
	handler := middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.RequestID,
		middleware.CORS,
		middleware.Logging,
		middleware.RateLimit(rateLimit("PRIVATE", privateConfigured, constants.PRIVATE_RATE_LIMIT), time.Minute),
	)

	// Configure server

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server (blocking)
	log.Printf("Server starting on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

}
