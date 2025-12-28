package main

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type clientLimiter struct {
	tokens     float64
	lastSeenAt time.Time
}

func newClientLimiter(burst int) *clientLimiter {
	return &clientLimiter{
		tokens:     float64(burst),
		lastSeenAt: time.Now(),
	}
}

func (cl *clientLimiter) allow(now time.Time, rps, burst int) bool {
	if rps <= 0 {
		return true
	}

	elapsed := now.Sub(cl.lastSeenAt).Seconds()
	if elapsed > 0 {
		cl.tokens += elapsed * float64(rps)
		if cl.tokens > float64(burst) {
			cl.tokens = float64(burst)
		}
	}

	if cl.tokens < 1 {
		cl.lastSeenAt = now
		return false
	}

	cl.tokens--
	cl.lastSeenAt = now
	return true
}

func (app *app) rateLimit(next http.Handler, rps, burst int) http.Handler {
	if next == nil {
		return nil
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*clientLimiter)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeenAt) > 5*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if ip == "" {
			app.respondError(w, r, http.StatusInternalServerError, "unable to determine client IP")
			return
		}
		now := time.Now()

		mu.Lock()
		limiter, ok := clients[ip]
		if !ok {
			limiter = newClientLimiter(burst)
			clients[ip] = limiter
		}
		canProceed := limiter.allow(now, rps, burst)
		mu.Unlock()

		if !canProceed {
			app.respondError(w, r, http.StatusTooManyRequests, "Too many requests please wait for sometime and try again")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	if real := r.Header.Get("X-Real-IP"); real != "" {
		return strings.TrimSpace(real)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
