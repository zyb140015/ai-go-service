# Project AI Rules

## 1. Project positioning

This is a **Go HTTP API service**.
Default code organization:

- `cmd`: process entrypoints and startup wiring only.
- `internal`: application code that must stay private to this module.
- `pkg`: reusable public helpers only when they are truly cross-project.
- `docs`: project documentation, operational notes, and AI rules.

If the real repository structure differs, **follow the repository instead of forcing a new layout**.

## 2. General principles

- New code uses Go and follows the current module structure.
- Prefer the **smallest correct change**.
- Do not invent commands, scripts, environment variables, package names, routes, tables, topics, or config keys.
- Read the existing code, config, and docs before editing.
- Do not add dependencies unless the current stack clearly cannot support the requirement.
- Do not mix unrelated refactors into a small fix.

## 3. Go code rules

- Keep APIs explicit and readable.
- Return wrapped errors with context using `%w`.
- Accept `context.Context` as the first parameter for request-scoped or I/O work.
- Keep interfaces small and define them near the consumer when possible.
- Avoid hidden global state and implicit side effects.
- Prefer standard library features before third-party abstractions.

## 4. HTTP layer boundaries

### Handlers must

- own HTTP concerns only: parsing requests, validating inputs, calling services, and shaping responses.
- return stable JSON response structures.
- map internal errors to client-safe messages.

### Handlers must not

- embed business workflows directly when that logic belongs in a service.
- access databases, external systems, or files directly without a dedicated lower layer.
- leak internal error details, stack traces, or secrets to clients.

## 5. Configuration rules

- Load configuration in one place.
- Give each config item a clear default or validation rule.
- Keep sensitive configuration out of logs.
- Do not scatter `os.Getenv` calls across the codebase.

## 6. Logging rules

- Use structured logging.
- Include enough context to trace the module, action, and key identifiers.
- Do not log passwords, tokens, API keys, cookies, or private payloads.
- Use `warn` for recoverable failures and `error` for unrecoverable failures.

## 7. Comment rules

- Exported packages, types, interfaces, functions, and constants must have Go-style comments.
- Comments must explain responsibility, constraints, edge cases, or intent.
- Complex business rules, concurrency control, retry behavior, fallback paths, and non-obvious validation need comments.
- Do not write noise comments that restate the code line by line.
- When code changes, update related comments in the same change so documentation stays trustworthy.

## 8. Data access rules

- Keep persistence logic behind a repository or store boundary.
- Keep SQL, migrations, and transaction management centralized.
- Do not couple HTTP handlers directly to database details.
- When schema names, table names, or queries are unknown, inspect the real source first.

## 9. Testing strategy

- Behavior changes require tests.
- Prefer table-driven tests when they improve coverage and readability.
- Use `httptest` for HTTP handler behavior.
- Keep unit tests fast and deterministic.
- Avoid fragile timing-based assertions when observable state can be asserted instead.

## 10. Tooling and validation

- Prefer repository-defined commands such as `make test`, `make vet`, and `make lint`.
- Validate the smallest scope that proves the change is correct.
- If validation cannot be run, say exactly what was not verified.

## 11. When AI edits code

1. Identify the affected layer first: `cmd / internal / pkg / docs / ci`.
2. Inspect existing patterns before creating new files or abstractions.
3. Reuse current helpers, response types, config loaders, and logging patterns.
4. For non-trivial work, list the minimal files affected.
5. Keep boundaries clear while editing.
6. Finish by updating tests, comments, and any required docs.
7. Report validation results and remaining risks.

## 12. Output preference

When explaining a code change, prefer this order:

1. Goal
2. Affected layers
3. Planned files
4. Key risks
5. Change details
6. Validation
7. Assumptions
