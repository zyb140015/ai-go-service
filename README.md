# ai-go-service

`ai-go-service` is a minimal Go HTTP API service template built for stable evolution.
It favors a small dependency surface, explicit boundaries, strong comments, and AI-friendly project rules.

## Stack

- Go `1.26`
- HTTP router: `chi`
- Database driver: `pgx/v5`
- Logging: standard library `log/slog`
- Testing: standard library `testing` + `httptest`
- Linting: `golangci-lint`
- SQL generation: `sqlc`

## Project layout

- `cmd/server`: application entrypoint
- `internal/app`: application wiring and lifecycle
- `internal/config`: environment-driven configuration
- `internal/http`: routing, handlers, and HTTP response helpers
- `internal/observability`: logging setup
- `internal/domain`: domain types boundary for future business logic
- `internal/store`: persistence boundary for repositories and PostgreSQL adapters
- `db/migrations`: SQL migrations
- `db/queries`: SQL files managed by `sqlc`
- `docs`: project and AI operation documents

## Getting started

```bash
go mod tidy
go run ./cmd/server
```

The service starts on `:8080` by default.
If `DATABASE_URL` is provided, the application opens a PostgreSQL pool during startup and checks it from `/readyz`.

## Endpoints

- `GET /healthz`
- `GET /readyz`

## Configuration

All configuration is loaded from environment variables.

- `HTTP_ADDR` default: `:8080`
- `HTTP_READ_TIMEOUT` default: `5s`
- `HTTP_READ_HEADER_TIMEOUT` default: `2s`
- `HTTP_WRITE_TIMEOUT` default: `10s`
- `HTTP_IDLE_TIMEOUT` default: `30s`
- `HTTP_SHUTDOWN_TIMEOUT` default: `10s`
- `LOG_LEVEL` default: `INFO`
- `DATABASE_URL` default: empty, which disables the PostgreSQL dependency

Duration values accept standard Go duration strings such as `5s` and `1m`.

## Development commands

```bash
make test
make vet
make build
make sqlc
```

If `golangci-lint` is installed locally, run:

```bash
make lint
```

## Comment standard

This project treats comments as part of the code contract.

- Exported packages, types, interfaces, functions, and constants require Go-style comments.
- Non-obvious control flow, boundary rules, and future extension points should be explained.
- Comments must explain intent and constraints, not restate obvious code.
- When behavior changes, related comments must be updated in the same change.

## Middleware and error behavior

- Every request is logged with method, path, status, bytes, duration, and request ID.
- Panics are recovered into a stable JSON error response.
- The API uses stable error codes such as `not_found`, `method_not_allowed`, `service_unavailable`, and `internal_error`.

## AI operation rules

See `AGENTS.md` and `docs/ai-rules.md` for the full AI execution rules used by this project.
