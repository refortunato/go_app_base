# Go Application Base

Base project template for building Go microservices following Clean Architecture + DDD principles.

## Using This Template

### Creating a New Project

To create a new project from this template:

```sh
./scripts/create-new-project.sh
```

Follow the interactive prompts to:
- Define your application name
- Configure the remote repository
- Automatically update all project references

### Removing Example Files

Once your project is set up, remove the example files and code:

```sh
./scripts/remove-examples.sh
```

This will clean up example entities, use cases, repositories, and controllers.

> **Important**: After using the scripts, you can safely delete the `./scripts` directory if you no longer need it.

## Project Structure

This project follows Clean Architecture and Domain-Driven Design (DDD) patterns:

- `cmd/server/` - Application entry point
- `configs/` - Configuration and database connection setup
- `internal/core/application/` - Use cases, repository interfaces, queries
- `internal/core/domain/` - Domain entities, value objects, domain services
- `internal/infra/` - Infrastructure implementations (repositories, web, dependencies)
- `internal/shared/` - Shared utilities across modules

### Architecture Styles

This project supports **two architectural styles** for different use cases:

#### 1. DDD with Clean Architecture (for complex business logic)
Used by modules: `example/`, `health/`

Structure:
```
internal/{module}/
  ├── core/
  │   ├── application/  # Use cases, repository interfaces
  │   └── domain/       # Entities, value objects, domain services
  └── infra/
      ├── repositories/ # Database implementations
      └── web/          # Controllers, routes
```

**Use when**: Complex business rules, multiple aggregates, rich domain logic

#### 2. 4-Tier Simplified Architecture (for CRUD operations)
Used by modules: `simple_module/`

Structure:
```
internal/{module}/
  ├── models/       # Data structures
  ├── repositories/ # Database access
  ├── services/     # Business logic
  ├── controllers/  # HTTP handlers
  ├── routes.go     # Route definitions
  └── module.go     # Dependency wiring
```

**Use when**: Simple CRUD, straightforward validations, minimal domain complexity

Both styles maintain **module independence** and use the same **dependency injection pattern**.

## Development Scripts

Automate common development tasks with the provided scripts:

### Creating New Modules
```sh
./scripts/create-module.sh
```
Creates a complete module structure with dependency wiring. Supports both DDD and 4-tier architectures.

📖 **[Complete Module Creation Guide](./docs/scripts/create-module-guide.md)**

### Creating New Entities
```sh
./scripts/create-entity.sh
```
Scaffolds complete CRUD operations for an entity within an existing module. Generates domain entities, repositories, use cases/services, controllers, and routes.

📖 **[Complete Entity Creation Guide](./docs/scripts/create-entity-guide.md)**

### Best Practices
- Use `create-module.sh` first to set up your module structure
- Then use `create-entity.sh` to add entities with full CRUD operations
- Choose DDD for complex business logic, 4-tier for simple CRUD
- All dependencies are automatically wired in the container

## Prerequisites

- Docker and Docker Compose
- Go 1.25.5+ (optional, for local development without Docker)

## How to Run

After cloning the repository, navigate to the project root and execute the commands below.

**Run in development mode with auto-reload**:
```sh
make dev
```

**Run production image (api)**:
```sh
make prod
```

Or if you prefer, you can run the commands directly with `docker-compose`:
```sh
# API only (default)
docker-compose up app-api

# API + Kafka
docker-compose --profile kafka up

# API + RabbitMQ
docker-compose --profile rabbitmq up

# All services
docker-compose --profile kafka --profile rabbitmq --profile grpc up

# Development with hot reload
docker-compose up app-dev
```

## Managing Dependencies

If you develop without having Go installed on your machine or prefer to ensure that the downloaded libraries are for the same version of Go that goes to production, you can run the command below which will organize your dependencies.

```sh
make go-mod-tidy
```

## Available Endpoints

### API Documentation (Swagger)
```
http://localhost:8080/swagger/index.html
```

Interactive API documentation with **Swagger UI**. Test endpoints directly in your browser.

**Generate/Update documentation**:
```sh
make swagger
```

📖 **[Complete Swagger Guide](./docs/implementation/swagger-guide.md)**

### Health Check
```http
GET /health
```

Returns `{"status": "OK"}` if the database connection is healthy.

### Example Resource
```http
GET /examples/:id
```

Returns an example resource by ID.

### Product Resource (Simple Module)
```http
GET    /products           # List all products (pagination: ?limit=10&offset=0)
GET    /products/:id       # Get product by ID
POST   /products           # Create new product
PUT    /products/:id       # Update product
DELETE /products/:id       # Delete product
```

Demonstrates a simpler 4-tier architecture for CRUD operations.

## Runtime Modes

This application can run in multiple modes depending on the first CLI argument:

- `api` (default when argument is blank)
- `kafka`
- `rabbitmq`
- `grpc`

In Kubernetes, prefer overriding the container args, for example:
- API pod: `args: ["api"]`
- Kafka consumer pod: `args: ["kafka"]`

## Environment Variables

All environment variables should be prefixed with `SERVER_APP_`. See the `.env.example` file for available configuration options.

## Architecture

This project follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer**: Pure business logic, entities, and domain services
- **Application Layer**: Use cases, repository interfaces, application services
- **Infrastructure Layer**: Database implementations, web framework, external services
- **Presentation Layer**: HTTP controllers, request/response DTOs
