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
- `internal/store/sqlcdb`: generated query code produced by `sqlc`
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
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- `POST /auth/change-password`
- `GET /auth/me`
- `GET /notes/`
- `GET /notes/{noteID}`
- `POST /notes/`
- `PUT /notes/{noteID}`
- `DELETE /notes/{noteID}`

### Auth API examples

Register a user:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","displayName":"Demo User","password":"password123"}'
```

Log in:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"password123"}'
```

Refresh tokens:

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H 'Content-Type: application/json' \
  -d '{"refreshToken":"<refresh-token>"}'
```

Log out:

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H 'Authorization: Bearer <access-token>'
```

Change password:

```bash
curl -X POST http://localhost:8080/auth/change-password \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access-token>' \
  -d '{"currentPassword":"password123","newPassword":"newpass123"}'
```

Read the current user:

```bash
curl http://localhost:8080/auth/me \
  -H 'Authorization: Bearer <token>'
```

### Notes API examples

Create a note:

```bash
curl -X POST http://localhost:8080/notes/ \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access-token>' \
  -d '{"title":"hello","body":"world"}'
```

List notes:

```bash
curl "http://localhost:8080/notes/?page=1&pageSize=10&sort=created_at&order=desc&q=hello" \
  -H 'Authorization: Bearer <access-token>'
```

Get one note:

```bash
curl http://localhost:8080/notes/1 \
  -H 'Authorization: Bearer <access-token>'
```

Update a note:

```bash
curl -X PUT http://localhost:8080/notes/1 \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <access-token>' \
  -d '{"title":"updated title","body":"updated body"}'
```

Delete a note:

```bash
curl -X DELETE http://localhost:8080/notes/1 \
  -H 'Authorization: Bearer <access-token>'
```

All `/notes` routes now require a valid bearer access token.

### Notes list query parameters

- `page`: 1-based page number, default `1`
- `pageSize`: items per page, default `20`, max `100`
- `q`: optional case-insensitive title filter
- `sort`: `created_at` or `title`, default `created_at`
- `order`: `asc` or `desc`, default `desc`

The list response includes a `meta` block with the resolved paging, sorting, filter, and total count.

## API docs

- Swagger UI: `http://localhost:8080/docs`
- OpenAPI spec: `http://localhost:8080/openapi.yaml`
- OpenAPI validation: `make openapi-lint`

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
- `AUTH_TOKEN_SECRET` default: empty, which disables auth handlers backed by real tokens
- `AUTH_TOKEN_TTL` default: `24h`

Duration values accept standard Go duration strings such as `5s` and `1m`.

## Database scaffold

- `db/migrations/000001_init.up.sql` creates the sample `app_notes` table.
- `db/queries/notes.sql` contains concrete CRUD queries for notes.
- `internal/store/postgres/note_repository.go` shows how to wrap generated `sqlc` code behind a repository boundary.

## Local PostgreSQL workflow

Start PostgreSQL locally:

```bash
make db-up
```

Set the database connection string:

```bash
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/ai_go_service?sslmode=disable'
```

Install `golang-migrate` locally:

```bash
brew install golang-migrate
```

Apply all migrations:

```bash
make migrate-up
```

Rollback the latest migration:

```bash
make migrate-down
```

Create a new migration file pair:

```bash
make migrate-create name=add_note_tags
```

If you want to apply only the current base SQL manually, use:

```bash
docker compose exec -T postgres psql -U postgres -d ai_go_service < db/migrations/000001_init.up.sql
```

If you need to revert the base SQL manually:

```bash
docker compose exec -T postgres psql -U postgres -d ai_go_service < db/migrations/000001_init.down.sql
```

Watch PostgreSQL logs:

```bash
make db-logs
```

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

Validate the OpenAPI document:

```bash
make openapi-lint
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
