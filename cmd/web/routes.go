package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	// Use the http.NewServeMux() function to initialize a new ServeMux (aka router)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// create file server that serves files out of "./ui/static/ directory"
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Pass the servemux as the 'next' parameter to the secureHeaders middlware
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else
	// Recall that the serveMux is also an http handler object
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
