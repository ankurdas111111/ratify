# Ratify API

Minimal Go REST API skeleton used to experiment with rate limiting, JSON responses and
hot-reload workflows.

## Stack

- Go 1.24+
- Standard library `net/http`
- [`github.com/tomasen/realip`](https://github.com/tomasen/realip) to capture the real client IP
- [`golang.org/x/time/rate`](https://pkg.go.dev/golang.org/x/time/rate) for per-IP token bucket limiting
- [`air`](https://github.com/air-verse/air) for live reload (optional but recommended)

## Features

- Clean `cmd/api` entry point with `app` struct encapsulating logging and routing
- JSON response helpers to keep handlers concise
- Central rate-limiter middleware that:
  - Buckets requests by client IP
  - Limits to `rps` requests/second with `burst` capacity
  - Evicts idle clients in the background
- Example handlers and routes under `/about` and `/test`
- Makefile targets for common tasks (`run`, `build`, `test`, `watch`, etc.)

## Project Layout

```
cmd/api/
  main.go        // server bootstrap
  routes.go      // route registration + handlers
  responses.go   // JSON helpers
  middleware.go  // IP-based rate limiter
.air.toml        // air watcher config
Makefile         // helper commands (watch, run, build, tidy)
```

## Getting Started

```bash
git clone https://github.com/<you>/ratify.git && cd ratify
cp .env.example .env   # if applicable
go mod tidy
```

### Development Modes

- `make run` – run the API once
- `make watch` or `make w` – run via Air with hot reload (requires `go install github.com/air-verse/air@latest`)
- `make build` – build the binary into `./bin/api`
- `make test` – run unit tests (placeholder for now)

### Manual Air Setup

```bash
go install github.com/air-verse/air@latest
export PATH="$PATH:$(go env GOPATH)/bin"   # ensure air is on PATH
make w
```

## API Endpoints

| Method | Path    | Description                     | Sample Response                          |
|--------|---------|---------------------------------|------------------------------------------|
| GET    | `/about`| Basic service metadata          | `{"name":"Gothic"}`                      |
| GET    | `/test` | Example health/test endpoint    | `{"name":"testing"}`                     |

Both endpoints are wrapped by the rate limiter (`2` req/sec with burst `4` by default). Excess
traffic receives `429 Too Many Requests` with a JSON error:

```json
{"error":"Too many requests please wait for sometime and try again"}
```

## Highlights

- Bootstrapped project structure, tooling, and live reload workflow
- Implemented custom token-bucket limiter with goroutine-based eviction
- Enforced consistent JSON responses and logging strategy

Feel free to fork and expand with real database access, authentication, or more routes. This
repository intentionally stays lean so reviewers can quickly understand my Go service patterns.

