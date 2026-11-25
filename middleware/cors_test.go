package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowsEverySite(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	req.Header.Set("Origin", "https://accesofugaz.com")

	CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("origin = %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("credentials = %q", rec.Header().Get("Access-Control-Allow-Credentials"))
	}
}

func TestAllowMethodsKeepsTheRouteVerbs(t *testing.T) {
	got := AllowMethods([]string{"POST", "GET", "POST", "PUT"})
	want := "GET, POST, PUT, OPTIONS"
	if got != want {
		t.Fatalf("methods = %q, want %q", got, want)
	}
}
