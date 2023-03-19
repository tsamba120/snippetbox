package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Handler (aka controller) for home endpoint
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Use ParseFiles() to read template files into template set. + error handling
	// File path must be relative to CWD or an absolute path
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...) // look into variadic functions
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Execute method of template set writes template content to the response body.
	// Last parameter for Execute is for any dynamic data we want to pass in. Nil for now
	if err = ts.Execute(w, nil); err != nil {
		app.serverError(w, err)
	}
}

// handler to show snippet
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...\n", id)
}

// handler to create a snippet
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet...\n"))
}
