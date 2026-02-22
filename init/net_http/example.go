package net_http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Config holds server configuration
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// App represents the application
type App struct {
	config *Config
	router *http.ServeMux
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

// customResponseWriter wraps http.ResponseWriter to capture status code
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Handle error responses
func (app *App) handleError(w http.ResponseWriter, r *http.Request, status int, message string, err interface{}) {
	log.Printf("Error: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Status:  status,
	}

	json.NewEncoder(w).Encode(response)
}

// JSON response helper
func (app *App) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			app.handleError(w, nil, http.StatusInternalServerError, "Failed to encode response", err)
		}
	}
}

// Routes
func (app *App) setupRoutes() {
	// Apply middleware to routes
	app.router.HandleFunc("/api/health", app.errorHandler(app.corsMiddleware(app.healthHandler)))
	app.router.HandleFunc("/api/users", app.errorHandler(app.corsMiddleware(app.usersHandler)))
	app.router.HandleFunc("/api/users/", app.errorHandler(app.corsMiddleware(app.userHandler)))

	// Public route without auth (example)
	app.router.HandleFunc("/", app.errorHandler(app.corsMiddleware(app.homeHandler)))
}

// Handlers
func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.handleError(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	response := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}

	app.jsonResponse(w, http.StatusOK, response)
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.handleError(w, r, http.StatusNotFound, "Page not found", nil)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Welcome to Go Server</h1><p>API endpoints available at /api/*</p>"))
}

func (app *App) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get all users
		users := []map[string]interface{}{
			{"id": 1, "name": "John Doe", "email": "john@example.com"},
			{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
		}
		app.jsonResponse(w, http.StatusOK, users)

	case http.MethodPost:
		// Create new user
		var user map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			app.handleError(w, r, http.StatusBadRequest, "Invalid request body", err)
			return
		}
		defer r.Body.Close()

		// Validate user data
		if name, ok := user["name"]; !ok || name == "" {
			app.handleError(w, r, http.StatusBadRequest, "Name is required", nil)
			return
		}

		user["id"] = 3 // In real app, this would be auto-generated
		app.jsonResponse(w, http.StatusCreated, user)

	default:
		app.handleError(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func (app *App) userHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path /api/users/{id}
	id := r.URL.Path[len("/api/users/"):]

	if id == "" {
		app.handleError(w, r, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get user by ID
		user := map[string]interface{}{
			"id":    id,
			"name":  "John Doe",
			"email": "john@example.com",
		}
		app.jsonResponse(w, http.StatusOK, user)

	case http.MethodPut:
		// Update user
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			app.handleError(w, r, http.StatusBadRequest, "Invalid request body", err)
			return
		}
		defer r.Body.Close()

		updates["id"] = id
		app.jsonResponse(w, http.StatusOK, updates)

	case http.MethodDelete:
		// Delete user
		app.jsonResponse(w, http.StatusNoContent, nil)

	default:
		app.handleError(w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

// Logging middleware
func (app *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		crw := &customResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(crw, r)

		// Log the request
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, crw.statusCode, time.Since(start))
	})
}

func main() {
	// Load configuration
	config := DefaultConfig()

	// Override with environment variables if needed
	if port := os.Getenv("PORT"); port != "" {
		config.Port = ":" + port
	}

	// Create and start the app
	app := NewApp(config)

	if err := app.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
