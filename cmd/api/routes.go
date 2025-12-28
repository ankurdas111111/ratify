package main

import "net/http"

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/about", app.getAbout())
	mux.Handle("/test", app.getTest())
	return app.rateLimit(mux, 1, 2)
}

func (app *app) getAbout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.respondJSON(w, r, http.StatusOK, map[string]string{"name": "Gothic"})
	})
}

func (app *app) getTest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.respondJSON(w, r, http.StatusOK, map[string]string{"name": "testing"})
	})
}
