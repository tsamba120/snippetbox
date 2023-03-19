package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	// Use the http.NewServeMux() function to initialize a new ServeMux (aka router)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// create file server that serves files out of "./ui/static/ directory"
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
