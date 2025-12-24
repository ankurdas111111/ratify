package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestApp() *app {
	return &app{
		logger: log.New(io.Discard, "", 0),
	}
}

func TestRespondJSON(t *testing.T) {
	application := newTestApp()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)

	payload := map[string]string{"message": "ok"}
	application.respondJSON(recorder, request, http.StatusCreated, payload)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d got %d", http.StatusCreated, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected content-type application/json got %s", recorder.Header().Get("Content-Type"))
	}

	body := strings.TrimSpace(recorder.Body.String())
	if !strings.Contains(body, `"message":"ok"`) {
		t.Fatalf("expected body to contain message json, got %s", body)
	}
}

func TestRespondError(t *testing.T) {
	application := newTestApp()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)

	application.respondError(recorder, request, http.StatusBadRequest, "boom")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, recorder.Code)
	}

	body := strings.TrimSpace(recorder.Body.String())
	if body != `{"error":"boom"}` {
		t.Fatalf("expected error json got %s", body)
	}
}
