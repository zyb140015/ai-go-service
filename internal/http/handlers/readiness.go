package handlers

import (
	"net/http"

	"ai-go-service/internal/http/response"
)

// ReadinessChecker reports whether the service can handle traffic safely.
type ReadinessChecker interface {
	Check(r *http.Request) error
}

// ReadyzHandler returns a readiness handler backed by dependency checks.
func ReadyzHandler(checker ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checker != nil {
			if err := checker.Check(r); err != nil {
				response.Error(w, http.StatusServiceUnavailable, response.CodeServiceUnavailable, "service is not ready")
				return
			}
		}

		response.JSON(w, http.StatusOK, HealthPayload{Status: "ready"})
	}
}
