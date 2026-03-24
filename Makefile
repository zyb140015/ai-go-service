GO ?= go
SQLC ?= sqlc

.PHONY: test vet build run lint tidy sqlc

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
