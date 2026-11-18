package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

// serverError() logs the error and sends Internal Server Error to the client.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError() sends stats error code response to the client.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// render() renders a page template.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// If Dev mode is active, do not cache.
	if app.devMode {
		newTs, err := buildTemplate("./ui/html/pages/" + page)
		if err != nil {
			app.serverError(w, r, err)
		}
		ts = newTs
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Write HTTP status code.
	w.WriteHeader(status)

	// Write contents of buffer.
	buf.WriteTo(w)
}

// newTemplateData() returns default templateData structure.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		TimeNow: time.Now(),
	}
}

// getEnv() returns environment variable if exist otherwise fallback string.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
