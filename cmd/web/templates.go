package main

import (
	"html/template"
	"path/filepath"

	"letsgo.rumi/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
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

		// Parse base template
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Parse allt the partial templates
		ts, err = ts.ParseGlob("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Parse page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template to the map
		cache[name] = ts
	}
	return cache, nil
}
