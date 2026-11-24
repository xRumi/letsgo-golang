package main

import (
	"html/template"
	"path/filepath"
	"time"

	"letsgo.rumi/internal/models"
)

type templateData struct {
	TimeNow  time.Time
	Snippet  models.Snippet
	Snippets []models.Snippet
	Form     SnippetCreateForm
	Flash    string
}

// humanDate() returns formated time.Time string.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Template functions.
var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache() returns a map containing parsed templates for all the
// template pages.
func newTemplateCache() (map[string]*template.Template, error) {
	// Initializd a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Get all filepaths that match the pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract filename
		name := filepath.Base(page)

		// Build template
		ts, err := buildTemplate(page)
		if err != nil {
			return nil, err
		}

		// Add the template to the map
		cache[name] = ts
	}
	return cache, nil
}

// buildTemplate() returns a parsed template for a template page.
func buildTemplate(page string) (*template.Template, error) {
	// Parse base template.
	ts, err := template.New(page).Funcs(templateFunctions).ParseFiles("./ui/html/base.tmpl")
	if err != nil {
		return nil, err
	}

	// Parse all the partial templates.
	ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Parse page template.
	ts, err = ts.ParseFiles(page)
	if err != nil {
		return nil, err
	}

	return ts, err
}
