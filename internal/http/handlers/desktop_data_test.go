package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-go-service/internal/http/handlers"
	"ai-go-service/internal/service"
)

func TestDesktopSystemStatsHandlerRejectsUnauthorized(t *testing.T) {
	t.Parallel()

	handler := handlers.DesktopSystemStatsHandler(service.NewDesktopDataService(nil))
	request := httptest.NewRequest(http.MethodGet, "/desktop/stats/system?range=week", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestDesktopSystemStatsHandlerReturnsData(t *testing.T) {
	t.Parallel()

	handler := handlers.DesktopSystemStatsHandler(service.NewDesktopDataService(nil))
	request := httptest.NewRequest(http.MethodGet, "/desktop/stats/system?range=week", nil)
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body service.DesktopSystemStatsResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.AlertCards) == 0 {
		t.Fatal("expected system stats alert cards")
	}
}

func TestDesktopUsageStatsHandlerReturnsData(t *testing.T) {
	t.Parallel()

	handler := handlers.DesktopUsageStatsHandler(service.NewDesktopDataService(nil))
	request := httptest.NewRequest(http.MethodGet, "/desktop/stats/usage?range=week&tenant=all", nil)
	request.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body service.DesktopUsageStatsResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Tenants) == 0 {
		t.Fatal("expected usage stats tenants")
	}
}
