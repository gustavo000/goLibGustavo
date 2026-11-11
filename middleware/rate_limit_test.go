package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitCountsTheVisitorAddress(t *testing.T) {
	t.Setenv("K_SERVICE", "")
	handler := RateLimit(2, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for i, port := range []string{"1001", "1002", "1003"} {
		req := httptestRequest("203.0.113.8:" + port)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		want := http.StatusNoContent
		if i == 2 {
			want = http.StatusTooManyRequests
		}
		if rec.Code != want {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, want)
		}
	}
}
