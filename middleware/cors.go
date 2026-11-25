package middleware

import (
	"net/http"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")
		next.ServeHTTP(w, r)
	})
}

// AllowMethods is the Access-Control-Allow-Methods value for one path.
// OPTIONS is included because the browser sends it before GET, POST, PUT, PATCH, or DELETE.
func AllowMethods(methods []string) string {
	present := map[string]bool{}
	var extra []string
	for _, method := range methods {
		method = strings.ToUpper(strings.TrimSpace(method))
		if method == "" || method == http.MethodOptions || present[method] {
			continue
		}
		present[method] = true
		switch method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			extra = append(extra, method)
		}
	}
	ordered := make([]string, 0, len(present)+1)
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		if present[method] {
			ordered = append(ordered, method)
		}
	}
	ordered = append(ordered, extra...)
	ordered = append(ordered, http.MethodOptions)
	return strings.Join(ordered, ", ")
}
