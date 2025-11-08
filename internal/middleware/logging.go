// Package middleware provides HTTP middleware functions for the application.
// This file contains the logging middleware.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter is a custom wrapper around http.ResponseWriter to capture the status code.
// This is necessary because the status code is not directly accessible from the ResponseWriter interface.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader overrides the original WriteHeader method to capture the status code.
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware logs the details of each incoming HTTP request.
// It logs the method, URI, protocol, status code, and duration of the request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request processing.
		start := time.Now()

		// Wrap the original response writer with our custom one.
		rw := &responseWriter{w, http.StatusOK}

		// Call the next handler in the chain.
		next.ServeHTTP(rw, r)

		// Calculate the duration of the request.
		duration := time.Since(start)

		// Log the formatted request details.
		log.Printf("[%s] %s %s (%d %s) - %s",
			r.Method,
			r.RequestURI,
			r.Proto,
			rw.status,
			http.StatusText(rw.status),
			duration)
	})
}
