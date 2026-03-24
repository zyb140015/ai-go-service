GO ?= go
SQLC ?= sqlc

.PHONY: test vet build run lint tidy sqlc db-up db-down db-logs

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
