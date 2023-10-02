package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/tsamba120/snippetbox/pkg/models"
)

// Handler (aka controller) for home endpoint
// NOTE: know difference between pointer receivers vs value receivers
// when implementing methods on a type
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	app.render(
		w,
		r,
		"home.page.tmpl",
		&templateData{
			Snippets: s,
		},
	)
}

// handler to show snippet
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(
		w,
		r,
		"show.page.tmpl",
		&templateData{
			Snippet: s,
		},
	)
}

// New createSnippetForm handler, which for now returns a placeholder reponse
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}

// handler to create a snippet
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// initialize a map to hold any validation errors
	errors := make(map[string]string)

	// validate title from form
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		// RuneCount gets number of chars, len() gets number of bytes
		errors["title"] = "This fied is too long (maximum is 100 characters)"
	}

	// validate content field
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	// validate expires
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	// if there are any errors, dump them in plain text HTTP response and return
	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	// pass data to SnippetModel.Insert() method, receiving the ID back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
