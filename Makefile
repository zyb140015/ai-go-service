GO ?= go
SQLC ?= sqlc
MIGRATE ?= migrate

.PHONY: test vet build run lint tidy sqlc db-up db-down db-logs migrate-up migrate-down migrate-create

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

build:
	$(GO) build ./cmd/server

run:
	$(GO) run ./cmd/server

lint:
	golangci-lint run

tidy:
	$(GO) mod tidy

sqlc:
	$(SQLC) generate

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-logs:
	docker compose logs -f postgres

migrate-up:
	$(MIGRATE) -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	$(MIGRATE) -path db/migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	$(MIGRATE) create -ext sql -dir db/migrations -seq $(name)
