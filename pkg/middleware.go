package pkg

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// LogRequest a small middleware function which logs request details
func LogRequest() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userAgent := r.Header.Get("User-Agent")
			log.Tracef(" ====> request [%s] path: [%s] [UA: %s]", r.Method, r.URL.Path, userAgent)
			next.ServeHTTP(w, r)
		})
	}
}
