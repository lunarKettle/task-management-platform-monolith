package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware логирует входящие HTTP-запросы
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		log.Printf("Completed %d %s in %v", ww.statusCode, http.StatusText(ww.statusCode), duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	log.Printf("WriteHeader called with status: %d", code)
	if rw.written {
		log.Printf("WriteHeader called again, skipping")
		return
	}
	rw.written = true
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
