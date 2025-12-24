package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *app) rateLimit(next http.Handler, rps, burst int) http.Handler {
	type client struct{
		limiter *rate.Limiter
		lastSeenAt time.Time
	}
	var(
		mu sync.Mutex 
		clients = make(map[string]*client)
	)

	go func() {
		for { 
			time.Sleep(time.Minute)
			mu.Lock()
		    for ip, client := range clients {
				if time.Since(client.lastSeenAt) > 5 * time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
     }()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)
		mu.Lock()
		defer mu.Unlock()
		if _, ok := clients[ip]; !ok {
			clients[ip] = &client{ 
				limiter: rate.NewLimiter(rate.Limit(rps), burst), 
				lastSeenAt: time.Now(),
			}
		}

		if !clients[ip].limiter.Allow() {
			app.respondError(w, r, http.StatusTooManyRequests, "Too many requests please wait for sometime and try again")
			return
		}

		next.ServeHTTP(w, r)
	})
}



