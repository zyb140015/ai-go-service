package handlers

import (
	"net/http"

	"ai-go-service/internal/http/apidocs"
)

// OpenAPIHandler serves the embedded OpenAPI specification.
func OpenAPIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	_, _ = w.Write(apidocs.OpenAPI)
}

// SwaggerUIHandler serves a lightweight Swagger UI page.
func SwaggerUIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(apidocs.SwaggerHTML)
}
