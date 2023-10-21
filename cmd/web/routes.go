package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our "standard" middleware
	// which will be used for every request
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Create new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later
	dynamicMiddleware := alice.New(app.session.Enable)

	// Make sure to register exact matches to endpoints first, before any wildcards
	// Note we are now able to use semantic urls
	// Instead of query strings, we'll use a full, descriptive path to the resource
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// create file server that serves files out of "./ui/static/ directory"
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Recall that the serveMux is also an http handler object

	// Return the 'standard' middleware chain followed by the servemux
	return standardMiddleware.Then(mux)
}
