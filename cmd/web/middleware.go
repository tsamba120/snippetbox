package main

import (
	"fmt"
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

// Middleware to log every HTTP request to the server
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
			next.ServeHTTP(w, r)
		},
	)
}

// Middleware to recover from a panic and provide meaningful error + message to client
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack)
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there was...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response
				// This header on the response acts as a trigger for the Go
				// server to automatically close the current connection
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}() // this parentheses calls the anonymous function
		next.ServeHTTP(w, r)
	})
}
