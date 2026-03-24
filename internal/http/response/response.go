package response

import (
	"encoding/json"
	"net/http"
)

// ErrorBody describes a stable JSON error shape for HTTP clients.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains machine-readable and human-readable error details.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// JSON writes a JSON response with the provided status code and payload.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

// Error writes a JSON error response using the shared error schema.
func Error(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, ErrorBody{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
