package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

type app struct {
	logger *log.Logger
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	application := &app{
		logger: logger,
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      application.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     logger,
	}

	logger.Printf("starting server on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}
