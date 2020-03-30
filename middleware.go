package sso

import (
	"log"
	"time"
	"net/http"

	"github.com/rs/xid"
)

func CORS(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*");
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,OPTIONS");
		if r.Method == "OPTIONS" { return; }
		f.ServeHTTP(w, r)
	})
}

type StatusCodeLogger struct {
	w http.ResponseWriter
	Code int
}

func (s StatusCodeLogger) Header() http.Header {
	return s.w.Header()
}

func (s StatusCodeLogger) Write(content []byte) (int, error) {
	return s.w.Write(content)
}

func (s StatusCodeLogger) WriteHeader(statusCode int) {
	s.Code = statusCode
	s.w.WriteHeader(statusCode)
}

func Logging(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := xid.New()
		rw := StatusCodeLogger{w, 200}
		log.Printf("[%s] Started [%s] %s", reqid.String(), r.Method, r.URL.Path)
		start := time.Now()
		f.ServeHTTP(rw, r)
		end := time.Now()
		diff := float64(end.Sub(start)) / float64(time.Microsecond)
		log.Printf("[%s] Completed [%s] %d  %s (%.2f μs)", reqid.String(), r.Method, rw.Code, r.URL.Path, diff)
	})
}

func ContentType(contentType string) func(http.Handler) http.Handler {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			f.ServeHTTP(w, r)
		})
	}
}

type ContentTypeDetector struct {
	w http.ResponseWriter
	contentType string
}

func (c ContentTypeDetector) Header() http.Header {
	return c.w.Header()
}

func (c ContentTypeDetector) Write(content []byte) (int, error) {
	if c.contentType == "" {
		c.contentType = http.DetectContentType(content)
		c.w.Header().Set("Content-Type", c.contentType)
	}
	return c.w.Write(content)
}

func (c ContentTypeDetector) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func DetectContentType(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := ContentTypeDetector{w, ""}
		// w.Header().Set("Content-Type", "application/json")
		f.ServeHTTP(c, r)
	})
}
