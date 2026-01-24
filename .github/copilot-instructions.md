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
      - `container.go`: orchestrates module initialization by calling each module's factory.
      - Returns a `Container` struct with all initialized modules.
      - Does NOT wire individual dependencies - delegates to module factories.
- `configs/`
  - Application configuration (`Conf` struct) and DB connection setup.
  - When adding new environment variables, ensure they are represented here.
  - `config.go`: defines `Conf` struct and `LoadConfig()` function.
  - `db_connection.go`: MySQL connection setup.
- `internal/`
  - Application code not intended for external import.
  - **Each directory within `internal/` represents a DDD module/bounded context** (e.g., `example/`, `health/`).
  - Each module is **completely independent** and self-contained with its own `core/` and `infra/` layers.
  - **Module structure** (e.g., `internal/example/`):
    - `core/` (application + domain layers)
      - `core/application/` (use cases, application services)
        - `usecases/`: orchestrates business flows. Each use case has its own file.
        - `services/`: application-level services (orchestration, cross-aggregate coordination).
        - `query/`: read/query handlers and DTOs (CQRS-style queries).
        - `repositories/`: repository interfaces used by use cases/queries.
      - `core/domain/` (pure domain)
        - `entities/`: domain entities. Each entity has its own file.
        - `services/`: domain services (when entity can't resolve logic or logic is repeated).
        - `valueobjects/`: value objects. Each value object has its own file.
    - `infra/` (infrastructure layer - adapters for domain needs)
      - `module.go`: **Module Factory** - wires all module dependencies (repos, use cases, controllers).
      - `repositories/`: concrete repository implementations (e.g., MySQL).
      - `web/`: web delivery layer
        - `controllers/`: HTTP controllers (receive dependencies via constructor).
        - `routes.go`: module-specific route registration (exports `RegisterRoutes(router, controller)`).
        - `routes.go`: module-specific route registration.
      - `module.go`: module factory that wires all dependencies (repositories, use cases, controllers).
  - `internal/shared/` (generic reusable code, framework-agnostic)
    - Shared utilities usable across multiple modules (keep it small and truly generic).
    - `logger/`: logger interface and implementation (slog).
    - `errors/`: application error handling.
    - `web/`: generic web infrastructure
      - `context/`: `WebContext` interface (framework-agnostic) and `GinContextAdapter`.
      - `server/`: `Server` interface, `GinServer` implementation, and factory with callback pattern.
      - `advisor/`: HTTP error response helpers (generic, framework-agnostic).
  - `internal/infra/` (shared infrastructure - cross-module)
    - `web/`: shared web components
      - `register_routes.go`: main route orchestrator (delegates to each module).

## Coding conventions (Go + Clean Arch + DDD)

### Language and style
- **All code must be written in English** (identifiers, errors, DTO names, etc.).
- Use the **minimum number of comments** possible; prefer clear names and small functions.
- Follow Go naming conventions:
  - Exported types/functions: **PascalCase**.
  - Unexported variables/functions (local scope): **camelCase** (Go standard; do not use snake_case).
  - Package names: short, lowercase.

### Domain rules
- If a struct has an **ID** and represents a business concept lifecycle, treat it as a **Domain Entity** and place it under `internal/{module}/core/domain/entities`.
- Use **Value Objects** for validated, immutable concepts under `internal/{module}/core/domain/valueobjects`.

### Use cases
- Every use case must expose a public method named `Execute`.
- Each use case must have **its own** input/output DTOs:
  - `Execute(input InputDto) (OutputDto, error)`
  - Keep DTOs in the use case package (or a dedicated `dto` subpackage) and keep them flat.
- Use cases must depend only on:
  - domain types from the same module
  - repository interfaces in `internal/{module}/core/application/repositories`
  - application services (if needed)
- Use cases are instantiated in `cmd/server/container/container.go`.

### Queries (read models)
- Query handlers should live in `internal/{module}/core/application/query`.
- Prefer returning query DTOs (read models) rather than domain entities when appropriate.

### Controllers and routing
- Controllers live in `internal/{module}/infra/web/controllers`.
- Controllers receive dependencies via constructor (from container).
- Controllers use `WebContext` interface from `internal/shared/web/context` (not Gin directly).
- **Each module registers its own routes** in `internal/{module}/infra/web/routes.go`.
- Each module exports a `RegisterRoutes(router *gin.Engine, controller)` function.
- The central orchestrator in `internal/infra/web/register_routes.go` calls each module's registration function.
- This keeps modules independent and prepares them for potential extraction into microservices.

### Repositories
- Repository interfaces must be defined in `internal/{module}/core/application/repositories`.
- Implementations must be in `internal/{module}/infra/repositories`.
- Inside each repository implementation:
  - Have a DB model struct (what is stored in MySQL).
  - Provide explicit mapping helpers named like `mapToDomain...` and `mapToDB...`.

### Dependency injection (Composition Root)
- **Container is in `cmd/server/container/container.go`** (Main Component in Clean Architecture).
- Container is a struct that holds references to all initialized modules.
- **Each module wires its own dependencies** via `NewModuleXYZ()` factory in `internal/{module}/infra/module.go`.
- `container.New(db, cfg)` creates the logger and initializes all modules.
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
      infraWeb.RegisterRoutes(c),  // passes container to route orchestrator
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
- Examples: logger interface/implementation, web server abstractions, context adapters, error handling, HTTP response advisors.

**`internal/infra/`**
- Contains **application-specific infrastructure adapters**.
- Implements interfaces defined in `core/` (repositories) or `shared/` (if generic contracts exist).
- Can depend on `shared/` and `core/`, but not vice versa.
- Examples: application-specific route definitions.

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
srv := server.NewGinServerWithRoutes(cfg.WebServerPort, infraWeb.RegisterRoutes(c))
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

## Creating a new module (step-by-step)

When creating a new module, follow this workflow to ensure proper bounded context isolation:

### 1. Create module directory structure
```bash
internal/
  └── {module_name}/
      ├── core/
      │   ├── application/
      │   │   ├── repositories/
      │   │   └── usecases/
      │   └── domain/
      │       ├── entities/
      │       ├── services/      (optional)
      │       └── valueobjects/  (optional)
      └── infra/
          ├── repositories/
          └── web/
              └── controllers/
```

### 2. Domain Layer (core/domain)
- **Entities**: Create in `internal/{module}/core/domain/entities/{entity_name}.go`
  - Include constructor `New{Entity}()` with validation
  - Include restore function `Restore{Entity}()` for reconstitution from DB
  - Keep private fields with public getters/setters
  - Add `Validate()` method

- **Value Objects** (optional): Create in `internal/{module}/core/domain/valueobjects/{vo_name}.go`
  - Immutable structs with validation
  - No identity, compared by value

### 3. Application Layer (core/application)
- **Repository Interfaces**: Create in `internal/{module}/core/application/repositories/{entity}_repository.go`
  - Define contract: `Save`, `FindById`, `Update`, `Delete`, etc.
  - Return domain entities, not DB models

- **Use Cases**: Create in `internal/{module}/core/application/usecases/{action}_{entity}.go`
  - Each use case = one file
  - Define input/output DTOs in the same file
  - Expose `Execute(input InputDTO) (OutputDTO, error)` method
  - Depend only on repository interfaces and domain entities

### 4. Infrastructure Layer (infra)
- **Repositories**: Create in `internal/{module}/infra/repositories/{entity}_mysql_repository.go`
  - Implement interface from `core/application/repositories`
  - Define DB model struct (private)
  - Provide `mapToDomain()` and `mapToDB()` functions

- **Controllers**: Create in `internal/{module}/infra/web/controllers/{entity}_controller.go`
  - Receive use cases via constructor
  - Use `WebContext` interface (not Gin directly)
  - Use `advisor` for error responses

- **Routes**: Create in `internal/{module}/infra/web/routes.go`
  ```go
  package web
  
  func RegisterRoutes(router *gin.Engine, controller *controllers.XController) {
      router.GET("/path/:id", func(ctx *gin.Context) {
          controller.Method(context.NewGinContextAdapter(ctx))
      })
  }
  ```

### 5. Module Factory
- Create `internal/{module}/infra/module.go`:
```go
package infra

type {Module}Module struct {
    {Entity}Controller *controllers.{Entity}Controller
    {Action}UseCase    *usecases.{Action}UseCase
    // Add more dependencies as needed
}

func New{Module}Module(db *sql.DB) *{Module}Module {
    // Wire repositories
    repo := repositories.New{Entity}MySQLRepository(db)
    
    // Wire use cases
    useCase := usecases.New{Action}UseCase(repo)
    
    // Wire controllers
    controller := controllers.New{Entity}Controller(*useCase)
    
    return &{Module}Module{
        {Entity}Controller: controller,
        {Action}UseCase:    useCase,
    }
}
```

### 6. Register in Container
- Update `cmd/server/container/container.go`:
```go
import moduleInfra "github.com/refortunato/go_app_base/internal/{module}/infra"

type Container struct {
    {Module}Module *moduleInfra.{Module}Module
    // ... other modules
}

func New(db *sql.DB, cfg *configs.Conf) (*Container, error) {
    // ...
    moduleModule := moduleInfra.New{Module}Module(db)
    
    return &Container{
        {Module}Module: moduleModule,
        // ...
    }
}
```

### 7. Register Routes
- Update `internal/infra/web/register_routes.go`:
```go
import moduleWeb "github.com/refortunato/go_app_base/internal/{module}/infra/web"

func RegisterRoutes(c *container.Container) func(*gin.Engine) {
    return func(router *gin.Engine) {
        moduleWeb.RegisterRoutes(router, c.{Module}Module.{Entity}Controller)
        // ... other modules
    }
}
```

### Module independence checklist
- ✅ Module does NOT import other modules (only `shared` and `configs`)
- ✅ All dependencies wired in `module.go`
- ✅ Routes registered in module's own `routes.go`
- ✅ Can be extracted to separate service without changes

## Additional Resources

For detailed implementation guides:
- **Routes Management**: See `docs/implementation/routes-management.md`
- **Dependency Management**: See `docs/implementation/dependency-management.md`
2. **Wait for approval** before applying code changes.
3. If the user approves and the plan remains the same, proceed without asking again for each step.
4. If the plan changes materially (new files, new behavior, new architecture), present an updated plan and ask again.
