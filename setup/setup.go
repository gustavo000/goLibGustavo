package setup

import (
	"fmt"
	"log"
	"net/http"
)

func InitServer() {
	// Define route handlers

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/register", helloHandler)
	http.HandleFunc("/health", healthHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Home page handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Welcome to the Go server!")
}

// Hello endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "* ")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
	w.WriteHeader(http.StatusOK)
	return
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "healthy", "service": "go-server"}`)
}
