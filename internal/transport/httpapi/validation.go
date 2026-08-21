package httpapi

import (
	"example.com/media-workflow/internal/domain"
	"net/http"
	"strings"
)

func RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(r.Header.Get("X-Tenant-ID")) == "" && r.URL.Path != "/healthz" && r.URL.Path != "/readyz" {
			writeError(w, domain.NewInvalid("X-Tenant-ID header required", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}
func LimitBody(next http.Handler, max int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, max)
		next.ServeHTTP(w, r)
	})
}
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Header.Get("Content-Type") != "application/json" {
			writeError(w, domain.NewInvalid("content-type must be application/json", nil))
			return
		}
		next.ServeHTTP(w, r)
	})
}
