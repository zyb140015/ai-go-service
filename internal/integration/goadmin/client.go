package goadmin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const bearerPrefix = "Bearer "

// Client wraps the HTTP communication with go-admin.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a go-admin HTTP client.
func NewClient(baseURL string, timeout time.Duration) *Client {
	trimmedBaseURL := strings.TrimRight(strings.TrimSpace(baseURL), "/")

	return &Client{
		baseURL: trimmedBaseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (client *Client) newRequest(ctx context.Context, method string, path string, body any, accessToken string) (*http.Request, error) {
	requestURL := client.baseURL + path

	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}

		requestBody = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if strings.TrimSpace(accessToken) != "" {
		req.Header.Set("Authorization", bearerPrefix+strings.TrimSpace(accessToken))
	}

	return req, nil
}

func decodeJSON[T any](response *http.Response, target *T) error {
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func businessMessage(msg string, fallback string) string {
	if strings.TrimSpace(msg) != "" {
		return strings.TrimSpace(msg)
	}

	return fallback
}
