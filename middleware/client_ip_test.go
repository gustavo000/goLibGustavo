package middleware

import (
	"net/http"
	"testing"
)

func TestClientIPUsesRemoteAddrLocally(t *testing.T) {
	t.Setenv("K_SERVICE", "")
	req := httptestRequest("203.0.113.8:54321")
	req.Header.Set("X-Forwarded-For", "198.51.100.10")
	if got := ClientIP(req); got != "203.0.113.8" {
		t.Fatalf("local address = %q", got)
	}
}

func TestClientIPUsesLastForwardedOnCloudRun(t *testing.T) {
	t.Setenv("K_SERVICE", "tickets-bff")
	req := httptestRequest("169.254.8.1:1234")
	req.Header.Set("X-Forwarded-For", "198.51.100.10, 203.0.113.9")
	if got := ClientIP(req); got != "203.0.113.9" {
		t.Fatalf("visitor address = %q", got)
	}
}

func TestClientIPFallsBackWhenForwardedIsEmpty(t *testing.T) {
	t.Setenv("K_SERVICE", "tickets-bff")
	req := httptestRequest("[2001:db8::1]:443")
	if got := ClientIP(req); got != "2001:db8::1" {
		t.Fatalf("fallback address = %q", got)
	}
}

func httptestRequest(remote string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://example.test/ready", nil)
	if err != nil {
		panic(err)
	}
	req.RemoteAddr = remote
	return req
}
