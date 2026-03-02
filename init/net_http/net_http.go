package net_http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gustavo000/goLibGustavo/middleware"
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

// NewApp creates a new application instance
func NewApp(config *Config) *App {
	if config == nil {
		config = DefaultConfig()
	}

	return &App{
		config: config,
		router: http.NewServeMux(),
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

// CORS middleware
func (app *App) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Change this in production!
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
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

// Start the server
func StartHttp(routes rest.Routes) {
	// Create router
	mux := http.NewServeMux()
	for _, rt := range routes {
		route := rt
		mux.HandleFunc(fmt.Sprintf("%s %s", route.Method, route.Pattern), func(w http.ResponseWriter, r *http.Request) {
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
	}

	// Apply middleware chain
	handler := middleware.Chain(
		mux,
		middleware.Recovery,
		middleware.RequestID,
		middleware.CORS,
		middleware.Logging,
		middleware.RateLimit(100, time.Minute),
		middleware.Auth,
	)

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server (blocking)
	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

}
