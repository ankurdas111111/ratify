package main

import (
	"net/http"

	"github.com/ankurdas111111/ratify/pkg/ratelimit"
)

func (app *app) rateLimit(next http.Handler, rps, burst int) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.respondError(w, r, http.StatusTooManyRequests, "Too many requests please wait for sometime and try again")
	})

	middleware := ratelimit.Middleware(rps, burst, handler)
	return middleware(next)
}
