package sso

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/xid"
)

// CORS headers for cross-domain requests
func (s *server) cors(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,OPTIONS")
		if r.Method == "OPTIONS" {
			return
		}
		f.ServeHTTP(w, r)
	})
}

// Capture the response code whenever it's written so it can be retrieved
type statusCodeLogger struct {
	w    http.ResponseWriter
	Code int
}

func (s statusCodeLogger) Header() http.Header {
	return s.w.Header()
}

func (s statusCodeLogger) Write(content []byte) (int, error) {
	return s.w.Write(content)
}

func (s statusCodeLogger) WriteHeader(statusCode int) {
	s.Code = statusCode
	s.w.WriteHeader(statusCode)
}

// Log details about each request
func (s *server) logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := xid.New()
		rw := statusCodeLogger{w, 200}
		log.Printf("[%s] Started [%s] %s", reqid.String(), r.Method, r.URL.Path)
		start := time.Now()
		f.ServeHTTP(rw, r)
		end := time.Now()
		diff := float64(end.Sub(start)) / float64(time.Microsecond)
		log.Printf("[%s] Completed [%s] %d  %s (%.2f μs)", reqid.String(), r.Method, rw.Code, r.URL.Path, diff)
	})
}

// Force a given content type
func (s *server) contentType(contentType string) func(http.Handler) http.Handler {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			f.ServeHTTP(w, r)
		})
	}
}
