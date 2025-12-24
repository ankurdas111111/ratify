package main

import (
	"encoding/json"
	"net/http"
)

// respondJSON writes a status code plus JSON payload. Headers were removed to keep it minimal.
func (app *app) respondJSON(w http.ResponseWriter, _ *http.Request, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		app.logger.Printf("json encode error: %v", err)
	}
}

func (app *app) respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	app.respondJSON(w, r, status, map[string]string{"error": message})
}
