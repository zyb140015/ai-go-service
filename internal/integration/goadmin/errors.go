package goadmin

import (
	"errors"
	"strings"
)

var (
	// ErrUnauthorized indicates the upstream rejected the caller credentials.
	ErrUnauthorized = errors.New("goadmin unauthorized")
	// ErrUpstreamFailure indicates the upstream service failed to complete the request.
	ErrUpstreamFailure = errors.New("goadmin upstream failure")
	// ErrInvalidResponse indicates the upstream responded with an unexpected payload.
	ErrInvalidResponse = errors.New("goadmin invalid response")
)

// BusinessError is returned when go-admin completes the request but reports failure in its JSON body.
type BusinessError struct {
	Code    int
	Message string
}

// Error implements error.
func (err *BusinessError) Error() string {
	if err == nil {
		return ""
	}

	if strings.TrimSpace(err.Message) == "" {
		return "goadmin request failed"
	}

	return err.Message
}
