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
    - `cmd/server/main.go`: thin entry point that boots config, DB, creates container, and starts a runtime mode (api/kafka/rabbitmq/grpc).
    - `cmd/server/container/`: **Composition Root** (dependency injection)
      - `container.go`: wires all application dependencies (repositories, use cases, controllers).
      - Returns a `Container` struct with all initialized dependencies.
      - This is the **only place** where objects are composed together.
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
  - `internal/infra/` (infrastructure layer - adapters for domain needs)
    - `repositories/`: concrete repository implementations (e.g., MySQL).
    - `config/`: infrastructure configuration details.
    - `web/`: web delivery layer
      - `controllers/`: HTTP controllers (receive container dependencies via constructor).
      - `routes/`: route registration (`routes.go` receives container and returns route setup function).
  - `internal/shared/` (generic reusable code, framework-agnostic)
    - Shared utilities usable across multiple layers/modules (keep it small and truly generic).
    - `logger/`: logger interface and implementation (slog).
    - `errors/`: application error handling.
    - `web/`: generic web infrastructure
      - `context/`: `WebContext` interface (framework-agnostic) and `GinContextAdapter`.
      - `server/`: `Server` interface, `GinServer` implementation, and factory with callback pattern.

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
- Use cases are instantiated in `cmd/server/container/container.go`.

### Queries (read models)
- Query handlers should live in `internal/core/application/query`.
- Prefer returning query DTOs (read models) rather than domain entities when appropriate.

### Controllers and routing
- Controllers live in `internal/infra/web/controllers`.
- Controllers receive dependencies via constructor (from container).
- Controllers use `WebContext` interface from `internal/shared/web/context` (not Gin directly).
- Routes are registered in `internal/infra/web/routes/routes.go`.
- `RegisterRoutes` function receives the container and returns a `RouteSetupFunc`.

### Repositories
- Repository interfaces must be defined in `internal/core/application/repositories`.
- Implementations must be in `internal/infra/repositories`.
- Inside each repository implementation:
  - Have a DB model struct (what is stored in MySQL).
  - Provide explicit mapping helpers named like `mapToDomain...` and `mapToDB...`.

### Dependency injection (Composition Root)
- **Container is in `cmd/server/container/container.go`** (Main Component in Clean Architecture).
- Container is a struct that holds all application dependencies.
- `container.New(db, cfg)` creates and wires all dependencies.
- **No global singletons** - dependencies are passed via constructor or container.
- `cmd/server/main.go` must remain thin: load config, create DB, create container, start server.

### Server creation
- Use factory pattern from `internal/shared/web/server/factory.go`.
- `NewGinServerWithRoutes(port, setupRoutes)` accepts a callback function.
- Example in main.go:
  ```go
  c := container.New(db, cfg)
  srv := server.NewGinServerWithRoutes(
      cfg.WebServerPort,
      routes.RegisterRoutes(c),  // passes container to routes
  )
  srv.Start()
  ```

### Environment variables
- When a new environment variable is needed:
  - Add it to `cmd/server/.env` and `cmd/server/.env.example`.
  - Always use the prefix `SERVER_APP_...`.
  - Update configuration mapping in `configs/config.go`.

## Architecture patterns and principles

### Separation of concerns (shared vs infra)

**`internal/shared/`**
- Contains **generic, reusable, domain-agnostic** code.
- Can be extracted to a separate library and used in multiple projects.
- Must NOT depend on `infra/` or `core/` packages.
- Examples: logger interface/implementation, web server abstractions, context adapters, error handling.

**`internal/infra/`**
- Contains **application-specific infrastructure adapters**.
- Implements interfaces defined in `core/` (repositories) or `shared/` (if generic contracts exist).
- Can depend on `shared/` and `core/`, but not vice versa.
- Examples: MySQL repositories, application-specific controllers, route definitions.

**Decision criteria:**
- If code has knowledge of **domain entities or business rules** → `infra/`
- If code is **technical and framework-agnostic** → `shared/`
- If code is **business logic** → `core/`

### Dependency Inversion Principle in action

The factory pattern in `shared/web/server/factory.go` uses **callback functions** to invert dependencies:

```go
// shared/web/server/factory.go (generic)
type RouteSetupFunc func(*gin.Engine)

func NewGinServerWithRoutes(port string, setupRoutes RouteSetupFunc) *GinServer {
    router := gin.Default()
    if setupRoutes != nil {
        setupRoutes(router)
    }
    return NewGinServer(router, port)
}

// infra/web/routes/routes.go (application-specific)
func RegisterRoutes(c *container.Container) func(*gin.Engine) {
    return func(router *gin.Engine) {
        router.GET("/health", func(ctx *gin.Context) {
            c.HealthController.HealthCheck(context.NewGinContextAdapter(ctx))
        })
    }
}

// cmd/server/main.go (composition)
c := container.New(db, cfg)
srv := server.NewGinServerWithRoutes(cfg.WebServerPort, routes.RegisterRoutes(c))
```

This ensures `shared/` doesn't know about `infra/`, respecting dependency direction.

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
