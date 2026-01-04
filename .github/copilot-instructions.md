# Copilot Instructions (go_app_base)

## Project context

This repository is a Go application template built to serve as a **base for other services**, following **Clean Architecture + DDD**.

- **Tech stack**
  - **Golang 1.25**
  - **MySQL 8.0+**
  - **Gin** web framework
- **Architecture**
  - Designed following **Clean Architecture + DDD**.
- **Baseline features**
  - `GET /health`: checks DB connectivity via a simple `SELECT 1`.
  - `Example` aggregate sample: domain entity + use case + repository + endpoint.

## Repository layout (directory map)

Use this map to decide where new code belongs. Prefer adding code in the correct layer instead of mixing concerns.

- `cmd/`
  - `cmd/server/`
    - Entry point of the application.
    - `cmd/server/main.go` boots config, DB, dependency injection, and starts a runtime mode (api/kafka/rabbitmq/grpc).
- `configs/`
  - Application configuration and DB connection setup.
  - When adding new environment variables, ensure they are represented here.
- `internal/`
  - Application code not intended for external import.
  - `internal/core/` (application + domain layers)
    - `internal/core/application/` (use cases, application services)
      - `usecases/`: orchestrates business flows.
      - `services/`: application-level services (orchestration, cross-aggregate coordination).
      - `query/`: read/query handlers and DTOs (CQRS-style queries).
      - `repositories/`: repository interfaces used by use cases/queries.
    - `internal/core/domain/` (pure domain)
      - `entities/`: domain entities.
      - `services/`: domain services (created when a entity could not resolve some domain logic by itself or when a code/method is repeated among entities).
      - `valueobjects/`: value objects.
  - `internal/infra/` (infrastructure layer)
    - `dependencies/`: dependency injection composition root (singletons). See `internal/infra/dependencies/dependencies.go`.
    - `repositories/`: concrete repository implementations (e.g., MySQL).
    - `config/`: infrastructure configuration details.
    - `web/`: web delivery layer
      - `controllers/`: HTTP controllers.
      - `webserver/`: Gin bootstrap/router definitions (see `internal/infra/web/webserver/gin_handler.go`).
      - `webcontext/`: request/response context adapters.
- `internal/shared/`
  - Shared utilities usable across multiple layers/modules (keep it small and truly generic).

## Coding conventions (Go + Clean Arch + DDD)

### Language and style
- **All code must be written in English** (identifiers, errors, DTO names, etc.).
- Use the **minimum number of comments** possible; prefer clear names and small functions.
- Follow Go naming conventions:
  - Exported types/functions: **PascalCase**.
  - Unexported variables/functions (local scope): **camelCase** (Go standard; do not use snake_case).
  - Package names: short, lowercase.

### Domain rules
- If a struct has an **ID** and represents a business concept lifecycle, treat it as a **Domain Entity** and place it under `internal/core/domain/entities`.
- Use **Value Objects** for validated, immutable concepts under `internal/core/domain/valueobjects`.

### Use cases
- Every use case must expose a public method named `Execute`.
- Each use case must have **its own** input/output DTOs:
  - `Execute(input InputDto) (OutputDto, error)`
  - Keep DTOs in the use case package (or a dedicated `dto` subpackage) and keep them flat.
- Use cases must depend only on:
  - domain types
  - repository interfaces in `internal/core/application/repositories`
  - application services (if needed)

### Queries (read models)
- Query handlers should live in `internal/core/application/query`.
- Prefer returning query DTOs (read models) rather than domain entities when appropriate.

### Controllers and routing
- Controllers live in `internal/infra/web/controllers`.
- Routes and Gin wiring live in `internal/infra/web/webserver/gin_handler.go`.
- Prefer REST conventions where possible (resources, proper verbs, status codes).

### Repositories
- Repository interfaces must be defined in `internal/core/application/repositories`.
- Implementations must be in `internal/infra/repositories`.
- Inside each repository implementation:
  - Have a DB model struct (what is stored in MySQL).
  - Provide explicit mapping helpers named like `mapToDomain...` and `mapToDB...`.

### Dependency injection
- Singleton composition and wiring live in `internal/infra/dependencies/dependencies.go`.
- `cmd/server/main.go` must remain a thin entry point: load config, create DB, call dependency initialization, start the selected runtime mode.

### Environment variables
- When a new environment variable is needed:
  - Add it to `cmd/server/.env` and `cmd/server/.env.example`.
  - Always use the prefix `SERVER_APP_...`.
  - Update configuration mapping in `configs/config.go`.

## Runtime modes (Kubernetes-friendly)

This application can run in multiple modes depending on the first CLI argument:

- `api` (default when argument is blank)
- `kafka`
- `rabbitmq`
- `grpc`

In Kubernetes, prefer overriding the container args, for example:

- API pod: `args: ["api"]`
- Kafka consumer pod: `args: ["kafka"]`

## Collaboration protocol (required)

Before making changes, the AI agent must:

1. Present an action plan (bulleted steps) and a short list of files it intends to touch.
2. **Wait for approval** before applying code changes.
3. If the user approves and the plan remains the same, proceed without asking again for each step.
4. If the plan changes materially (new files, new behavior, new architecture), present an updated plan and ask again.
