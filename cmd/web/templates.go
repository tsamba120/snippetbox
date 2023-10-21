package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/tsamba120/snippetbox/pkg/forms"
	"github.com/tsamba120/snippetbox/pkg/models"
)

// TemplateData struct to act as holding structure for any
// dynamic data that we pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	CurrentYear int
	Flash		string
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Function which returns formatted tiemstamp string
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Essentially a key-value store for template functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize new map which acts as a cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all functions with
	// the extenstion '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the applications
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one
	for _, page := range pages {
		// Extract the file name like ('home.page.tmpl') from the full file path
		// and assign it to the name variable
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// momment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add template set to in-mem cache
		cache[name] = ts

	}
	return cache, nil
}
