# Go Application Base

Base project template for building Go microservices following Clean Architecture + DDD principles.

## Project Structure

This project follows Clean Architecture and Domain-Driven Design (DDD) patterns:

- `cmd/server/` - Application entry point
- `configs/` - Configuration and database connection setup
- `internal/core/application/` - Use cases, repository interfaces, queries
- `internal/core/domain/` - Domain entities, value objects, domain services
- `internal/infra/` - Infrastructure implementations (repositories, web, dependencies)
- `internal/shared/` - Shared utilities across modules

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
