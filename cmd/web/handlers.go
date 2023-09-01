package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tsamba120/snippetbox/pkg/models"
)

// Handler (aka controller) for home endpoint
// NOTE: know difference between pointer receivers vs value receivers
// when implementing methods on a type
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

// handler to create a snippet
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// for now we'll create some variables holding mock data. remove later on
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n-Kobayashi Issa"
	expires := "7"

	// pass data to SnippetModel.Insert() method, receiving the ID back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
