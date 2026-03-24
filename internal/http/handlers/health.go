package handlers

import (
	"net/http"

	"ai-go-service/internal/http/response"
)

// HealthPayload is returned by health and readiness endpoints.
type HealthPayload struct {
	Status string `json:"status"`
}

// Healthz reports whether the process is alive.
func Healthz(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, HealthPayload{Status: "ok"})
}
