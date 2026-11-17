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
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

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

func buildTemplate(page string) (*template.Template, error) {
	// Parse base template
	ts, err := template.New(page).Funcs(templateFunctions).ParseFiles("./ui/html/base.tmpl")
	if err != nil {
		return nil, err
	}

	// Parse allt the partial templates
	ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Parse page template
	ts, err = ts.ParseFiles(page)
	if err != nil {
		return nil, err
	}

	return ts, err
}
