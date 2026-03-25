package handlers_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "ai-go-service/internal/http"
)

type failingChecker struct{}

func (failingChecker) Check(_ *http.Request) error {
	return assertErr{}
}

type assertErr struct{}

func (assertErr) Error() string {
	return "dependency unavailable"
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHealthz(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	httpserver.NewRouter(newTestLogger(), nil, nil, nil, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status %q, got %q", "ok", body["status"])
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	recorder := httptest.NewRecorder()

	httpserver.NewRouter(newTestLogger(), nil, nil, nil, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Error.Code != "not_found" {
		t.Fatalf("expected error code %q, got %q", "not_found", body.Error.Code)
	}
}

func TestReadyzReturnsServiceUnavailableWhenDependencyFails(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()

	httpserver.NewRouter(newTestLogger(), failingChecker{}, nil, nil, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, recorder.Code)
	}

	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Error.Code != "service_unavailable" {
		t.Fatalf("expected error code %q, got %q", "service_unavailable", body.Error.Code)
	}
}
