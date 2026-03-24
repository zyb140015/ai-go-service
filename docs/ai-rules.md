# AI Rules For `ai-go-service`

## Goal

This document tells AI agents how to work in this repository safely.
The priority is correctness, small diffs, clear boundaries, and maintainable comments.

## Project shape

- `cmd/server`: starts the process and owns OS lifecycle integration.
- `internal/app`: wires configuration, logger, router, and future services.
- `internal/config`: loads and validates environment-driven settings.
- `internal/http`: owns transport-layer concerns.
- `internal/observability`: shared logging and future tracing setup.
- `internal/domain`: future business types and use-case boundaries.
- `internal/store`: future persistence adapters and repositories.

## Required execution order

1. Identify the affected layer.
2. Inspect existing code and naming patterns.
3. Reuse current helpers and structures.
4. Make the smallest correct change.
5. Update tests, comments, and docs if behavior changes.
6. Run the smallest meaningful validation.
7. Report risks and assumptions.

## Editing rules

- Do not introduce speculative abstractions.
- Do not rewrite package layout without a real need.
- Do not rename routes, config keys, or response fields casually.
- Do not add dependencies unless there is a clear gap in the current stack.
- Keep comments aligned with behavior.

## Comment quality standard

Good comments are required in this repository.

- Document exported symbols in Go style.
- Explain why code exists, what boundary it protects, or what edge case it handles.
- Add comments to non-obvious validation, concurrency behavior, shutdown logic, retries, or fallbacks.
- Avoid comments that simply narrate obvious statements.
- Remove or update stale comments immediately when behavior changes.

## HTTP boundary rules

- Parse, validate, and map requests in handlers.
- Keep business logic out of handlers when it starts to grow.
- Use shared response types for stable client contracts.
- Return safe error messages to clients and richer context to logs.

## Configuration rules

- Read configuration centrally.
- Validate invalid or empty values early.
- Do not scatter environment access across the codebase.

## Logging rules

- Use structured logs.
- Include enough context to debug production issues.
- Never log secrets or full sensitive payloads.

## Testing rules

- Add or update tests for behavior changes.
- Prefer `httptest` for HTTP behavior.
- Prefer fast deterministic tests.

## Validation commands

```bash
make test
make vet
make build
make lint
```
