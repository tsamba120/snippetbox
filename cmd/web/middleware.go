package main

import (
	"net/http"
)

// Middleware pattern using an anonymous function to wrap
// main request functionality
// This middleware will act on every request, so it will be
// executed before a request hits the servemux
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Frame-Options", "deny")
			next.ServeHTTP(w, r)
		})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
			next.ServeHTTP(w, r)
		},
	)
}
