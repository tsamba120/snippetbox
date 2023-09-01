package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our "standard" middleware
	// which will be used for every request
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Use the http.NewServeMux() function to initialize a new ServeMux (aka router)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// create file server that serves files out of "./ui/static/ directory"
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Recall that the serveMux is also an http handler object

	// Return the 'standard' middleware chain followed by the servemux
	return standardMiddleware.Then(mux)
}
