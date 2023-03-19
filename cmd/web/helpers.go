package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
