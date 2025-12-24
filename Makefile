.PHONY: run build test watch w tidy clean

GOFILES := ./cmd/api
BIN := ./bin/api

run:
	GO111MODULE=on go run $(GOFILES)

build:
	GO111MODULE=on go build -o $(BIN) $(GOFILES)

test:
	go test ./...

watch:
	air

w: watch

tidy:
	go mod tidy

clean:
	rm -rf $(BIN) tmp .gocache

