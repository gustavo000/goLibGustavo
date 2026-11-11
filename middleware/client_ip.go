package middleware

import (
	"net"
	"net/http"
	"os"
	"strings"
)

// ClientIP is the visitor address used by the rate limiter.
// Cloud Run sets K_SERVICE and appends the visitor to X-Forwarded-For, so the
// last value is the address Google observed. Local runs keep RemoteAddr.
func ClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if os.Getenv("K_SERVICE") != "" {
		if ip := lastForwarded(r.Header.Get("X-Forwarded-For")); ip != "" {
			return ip
		}
	}
	return hostOnly(r.RemoteAddr)
}

func lastForwarded(header string) string {
	parts := strings.Split(header, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		if ip := hostOnly(parts[i]); ip != "" {
			return ip
		}
	}
	return ""
}

func hostOnly(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}
