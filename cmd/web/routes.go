package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes() sets all the routes for the server.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Routes.
	mux.Handle("GET /static/", http.FileServer(http.Dir("./ui/")))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)
	mux.HandleFunc("GET /snippet/latest", app.snippetLatest)

	// Middleware chain.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
