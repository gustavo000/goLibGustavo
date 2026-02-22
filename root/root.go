package root

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gustavo000/goLibGustavo/init/external_services"
	"github.com/gustavo000/goLibGustavo/init/net_http"
	"github.com/gustavo000/goLibGustavo/models/properties"
	"github.com/gustavo000/goLibGustavo/models/rest"
	"github.com/gustavo000/goLibGustavo/routing"
)

var service = []*rest.Service{
	{
		Name:    "Core",
		Ingress: "http://cust-rtmn-orch-dl-core-{release}-service.{namespace}.svc.cluster.local",
		Path:    "/dl-core",
		Timeout: 30,
	},
}

func InitServer(routes rest.Routes) error {
	properties.NewProperties(
		properties.WithServiceName("github.com/gustavo000/goLibGustavo"),
		properties.WithBasePath("/library"),
		properties.WithPort("8080"),
		properties.WithNameSpace("orch"),
		properties.WithRelease("r1"),
		properties.WithVersion(getCurrentVersion()),
		properties.WithEnvironment("PROD"),
		properties.WithServices(service...))
	external_services.GetEndpointsOfServices(properties.GetProperty().GetEnv())
	app := net_http.StartHttp(routes)

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Error().Msg(fmt.Sprintf("failed to shutdown TracerProvider: %v", err.Error()))
		}
	}()

	errListen := app.Listen(":" + properties.GetProperty().GetPort())
	return errListen
}

func getCurrentVersion() string {
	var res map[string]string
	bytes, _ := os.ReadFile("tag_version.json")
	err := json.Unmarshal(bytes, &res)
	if err != nil {
		return ""
	}
	return res["version"]
}

func InitServer1() {
	// Define route handlers

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/register", helloHandler)
	http.HandleFunc("/health", healthHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}

func StartServer() {

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
