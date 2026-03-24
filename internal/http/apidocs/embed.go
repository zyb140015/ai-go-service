package apidocs

import _ "embed"

// OpenAPI contains the embedded OpenAPI specification served by the HTTP layer.
//
//go:embed openapi.yaml
var OpenAPI []byte

// SwaggerHTML contains a lightweight Swagger UI page that loads the embedded spec.
//
//go:embed swagger.html
var SwaggerHTML []byte
