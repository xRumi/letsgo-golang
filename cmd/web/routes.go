package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes() sets all the routes for the server.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Static files route. We don't want session middleware here.
	mux.Handle("GET /static/", http.FileServer(http.Dir("./ui/")))

	// Session middleware.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/latest", dynamic.ThenFunc(app.snippetLatest))

	// Middleware chain.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
