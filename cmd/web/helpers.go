package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// write error message to app errorLog
// both appends to error log and returns http error response to client
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// debug.Stack() gets stack trace for current goroutine and appends to log
	app.errorLog.Output(2, trace) // change stack depth so logger logs where error actually occurs! not where log is written

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// clientError sends specific status code and description to user from client
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// wrapper around clientError specificially for 404 not found errors
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// addDefaultData helper that takes a pointer to a templateDat struct
// adds the current year to the CurrentYear field, and then returns
// the pointer. We are not using the *http.Request parameter at the moment
// but we'll add it later
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	return td
}

func (app *application) render(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	td *templateData,
) {
	// Retrieve the appropriate template set from the cached base on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize buffer for two-stage render
	// If writing template to buffer fails, we send user an error message
	// Otherwise, we can write contents in buffer to the http.ResponseWriter
	buf := new(bytes.Buffer)

	// Execute the template set, passing in any dynamic data
	// Write to our buffer. If there is an error throw a serverError
	td = app.addDefaultData(td, r)
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
	}

	// Write contents of buffer to http.ResponseWriter
	buf.WriteTo(w)
}
