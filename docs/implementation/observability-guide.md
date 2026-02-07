# Observability Guide - OpenTelemetry + Jaeger

This guide explains how observability is implemented in this project using **OpenTelemetry** for instrumentation and **Jaeger** for distributed tracing visualization.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Auto-Instrumentation](#auto-instrumentation)
- [Custom Instrumentation](#custom-instrumentation)
- [Accessing Jaeger UI](#accessing-jaeger-ui)
- [Understanding Traces](#understanding-traces)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Production Considerations](#production-considerations)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                       Application                            │
│  ┌──────────────┐   ┌────────────┐   ┌──────────────┐      │
│  │ HTTP Request │──▶│  Use Case  │──▶│  Repository  │      │
│  │ (Gin)        │   │  (Domain)  │   │  (Database)  │      │
│  └──────┬───────┘   └─────┬──────┘   └──────┬───────┘      │
│         │                 │                  │               │
│         │ Span 1          │ Span 2           │ Span 3        │
│         ▼                 ▼                  ▼               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │       OpenTelemetry SDK (Tracer Provider)            │   │
│  └───────────────────────────┬──────────────────────────┘   │
└──────────────────────────────┼──────────────────────────────┘
                               │
                               ▼
                     ┌─────────────────────┐
                     │  Jaeger Collector   │
                     │  (OTLP HTTP: 4318)  │
                     └──────────┬──────────┘
                                │
                                ▼
                     ┌─────────────────────┐
                     │   Jaeger Storage    │
                     │   (In-Memory/Dev)   │
                     └──────────┬──────────┘
                                │
                                ▼
                     ┌─────────────────────┐
                     │    Jaeger UI        │
                     │  (Port 16686)       │
                     └─────────────────────┘
```

**Key Components:**

1. **OpenTelemetry SDK**: Instruments the application code
2. **Tracer Provider**: Manages trace lifecycle and exporters
3. **OTLP Exporter**: Sends traces to Jaeger via HTTP (port 4318)
4. **Jaeger**: Collects, stores, and visualizes traces

---

## Quick Start

### 1. Start Environment

```bash
make dev
```

This starts:
- MySQL (port 3306)
- Application (port 8080)
- Jaeger UI (port 16686)

### 2. Generate Traffic

```bash
# Health check (no tracing)
curl http://localhost:8080/health

# Example endpoint (with tracing)
curl http://localhost:8080/examples/550e8400-e29b-41d4-a716-446655440000
```

### 3. View Traces

```bash
make jaeger-ui
# Or manually: http://localhost:16686
```

---

## Configuration

### Environment Variables

Add to `cmd/server/.env`:

```env
# Enable/disable OpenTelemetry tracing
SERVER_APP_OTEL_ENABLED=true

# Service name in traces
SERVER_APP_OTEL_SERVICE_NAME=go_app_base

# Jaeger collector endpoint (Docker service name)
SERVER_APP_JAEGER_ENDPOINT=jaeger:4318
```

### Configuration Struct

In [`configs/config.go`](../../configs/config.go):

```go
type Conf struct {
    // ... other fields
    OtelEnabled     bool   `mapstructure:"SERVER_APP_OTEL_ENABLED"`
    OtelServiceName string `mapstructure:"SERVER_APP_OTEL_SERVICE_NAME"`
    JaegerEndpoint  string `mapstructure:"SERVER_APP_JAEGER_ENDPOINT"`
}
```

### Docker Networking

**Critical**: Use **container service names**, not `localhost`:

```yaml
# ✅ Correct (Docker network)
SERVER_APP_JAEGER_ENDPOINT=jaeger:4318

# ❌ Wrong (localhost doesn't work in Docker)
SERVER_APP_JAEGER_ENDPOINT=localhost:4318
```

---

## Auto-Instrumentation

### HTTP Requests (Gin Middleware)

**Location**: [`internal/shared/observability/middleware.go`](../../internal/shared/observability/middleware.go)

**What it does:**
- Creates a span for every HTTP request
- Captures: method, path, status code, duration
- Propagates trace context (W3C Trace Context headers)

**Applied in**: [`internal/shared/web/server/factory.go`](../../internal/shared/web/server/factory.go)

```go
if otelEnabled {
    router.Use(observability.TracingMiddleware(serviceName))
}
```

**Span attributes captured:**
- `http.method` (GET, POST, etc.)
- `http.url` (request path)
- `http.status_code` (200, 404, 500, etc.)
- `http.user_agent` (client info)

### Database Queries (SQL Wrapper)

**Location**: [`internal/shared/observability/db_tracer.go`](../../internal/shared/observability/db_tracer.go)

**What it does:**
- Wraps `*sql.DB` to add tracing to all queries
- Captures: SQL statement, duration, errors

**Applied in**: [`cmd/server/container/container.go`](../../cmd/server/container/container.go)

```go
if cfg.OtelEnabled {
    tracedDB = observability.WrapDB(db)
}
```

**Span attributes captured:**
- `db.system` (mysql)
- `db.statement` (SQL query)
- `db.operation` (query, exec, query_row)

---

## Custom Instrumentation

### Use Cases (Business Logic)

Add custom spans to use cases for better visibility into business operations.

**Example**: [`internal/example/core/application/usecases/get_example.go`](../../internal/example/core/application/usecases/get_example.go)

```go
import (
    "context"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
)

func (u *GetExampleUseCase) Execute(ctx context.Context, input GetExampleInputDTO) (*GetExampleOutputDTO, error) {
    // Create a span for this use case
    tracer := otel.Tracer("example.usecase")
    ctx, span := tracer.Start(ctx, "GetExampleUseCase.Execute")
    defer span.End()

    // Add custom attributes
    span.SetAttributes(
        attribute.String("example.id", input.Id),
        attribute.String("usecase", "GetExample"),
    )

    // Business logic...
    example, err := u.exampleRepository.FindById(input.Id)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "Failed to find example")
        return nil, err
    }

    span.SetStatus(codes.Ok, "Example retrieved successfully")
    return output, nil
}
```

**Key points:**
1. Always accept `context.Context` as first parameter
2. Use `tracer.Start(ctx, "span-name")` to create span
3. Always `defer span.End()`
4. Add relevant attributes for debugging
5. Record errors with `span.RecordError(err)`
6. Set status with `span.SetStatus(codes.Ok/Error, message)`

### Controllers (Context Propagation)

Controllers must extract context from HTTP request and pass it to use cases.

**Example**: [`internal/example/infra/web/controllers/example_controller.go`](../../internal/example/infra/web/controllers/example_controller.go)

```go
func (controller *ExampleController) GetExample(c webcontext.WebContext) {
    // Extract context from HTTP request (contains trace info)
    ctx := c.GetContext()

    // Pass context to use case
    output, err := controller.GetExampleUseCase.Execute(ctx, input)
    // ...
}
```

---

## Accessing Jaeger UI

### Option 1: Make Command

```bash
make jaeger-ui
```

### Option 2: Manual

Open browser to: **http://localhost:16686**

### Jaeger UI Features

1. **Search Traces**:
   - Service: `go_app_base`
   - Operation: `GET /examples/:id`
   - Tags: `http.status_code=200`

2. **Trace Timeline**:
   - See full request flow
   - Parent-child span relationships
   - Duration of each operation

3. **Span Details**:
   - Tags (attributes)
   - Logs (errors, events)
   - Process info (service name, version)

---

## Understanding Traces

### Trace Hierarchy

```
Trace (complete request flow)
 └── HTTP Span (Gin middleware)
      ├── Use Case Span (business logic)
      │    └── Repository Span (database query)
      │         └── SQL Span (actual query execution)
      └── Response
```

### Example Trace

**Request**: `GET /examples/550e8400-e29b-41d4-a716-446655440000`

**Spans**:
1. **HTTP Request** (`otelgin.Middleware`)
   - Duration: 45ms
   - Tags: `http.method=GET`, `http.status_code=200`

2. **GetExampleUseCase.Execute** (use case)
   - Duration: 42ms
   - Tags: `example.id=550e8400...`, `usecase=GetExample`

3. **sql.query** (database)
   - Duration: 40ms
   - Tags: `db.statement=SELECT ...`, `db.system=mysql`

**Insights**:
- Most time spent in database (40ms / 45ms = 89%)
- HTTP overhead: 3ms
- Use case overhead: 2ms

---

## Best Practices

### 1. When to Add Custom Spans

✅ **DO add spans for:**
- Use case entry points (`Execute` methods)
- Complex business operations (multi-step workflows)
- External service calls (HTTP, gRPC, message queues)
- Heavy computations (data processing, validations)

❌ **DON'T add spans for:**
- Simple getters/setters
- Trivial functions (< 1ms execution time)
- Every function (creates noise)

### 2. Naming Conventions

**Format**: `{layer}.{operation}`

Examples:
- `usecase.GetExample` ✅
- `repository.SaveExample` ✅
- `service.CalculateTotal` ✅
- `http.GET /examples/:id` (auto-generated) ✅

### 3. Attributes (Tags)

**Always include:**
- Business identifiers (IDs, user IDs, order numbers)
- Operation type (create, update, delete)
- Important parameters

**Example:**
```go
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.String("order.id", orderID),
    attribute.String("operation", "create_order"),
    attribute.Int("items.count", len(items)),
)
```

### 4. Error Handling

**Always record errors:**
```go
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, "descriptive error message")
    return err
}
```

### 5. Context Propagation

**Golden rule**: Always accept and pass `context.Context`

```go
// ✅ Good
func (u *UseCase) Execute(ctx context.Context, input DTO) (DTO, error) {
    // ...
}

// ❌ Bad (loses trace context)
func (u *UseCase) Execute(input DTO) (DTO, error) {
    // ...
}
```

---

## Troubleshooting

### Problem: No traces in Jaeger

**Check:**
1. Is `SERVER_APP_OTEL_ENABLED=true`?
2. Is Jaeger container running? (`docker ps | grep jaeger`)
3. Is application connecting to Jaeger? (Check logs for "OpenTelemetry tracing initialized")
4. Did you generate traffic? (No requests = no traces)

### Problem: "connection refused" error

**Cause**: Application can't reach Jaeger collector

**Solution**: Verify Docker networking
```bash
# Check if jaeger container is running
docker ps | grep jaeger

# Check if app can reach jaeger
docker exec -it go_app_base_dev ping jaeger

# Verify environment variable
docker exec -it go_app_base_dev env | grep JAEGER
```

### Problem: Traces not linked (parent-child relationship broken)

**Cause**: Context not propagated correctly

**Solution**:
1. Ensure use case methods accept `context.Context`
2. Extract context from HTTP request: `ctx := c.GetContext()`
3. Pass context to repository methods

### Problem: High memory/CPU usage

**Cause**: Too many spans or always sampling

**Solution**:
1. Reduce sampling rate in production (currently 100%)
2. Remove unnecessary spans
3. Use batch exporter (already configured)

---

## Production Considerations

### 1. Sampling Strategy

**Development** (current):
```go
sdktrace.WithSampler(sdktrace.AlwaysSample())
```

**Production** (recommended):
```go
// Sample 10% of traces
sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1))
```

### 2. Storage Backend

**Development** (current):
- In-memory (data lost on restart)

**Production** (recommended):
- Elasticsearch
- Cassandra
- Cloud providers (AWS X-Ray, Google Cloud Trace)

### 3. Exporter Configuration

**Current**: OTLP HTTP to Jaeger

**Production alternatives**:
- **AWS X-Ray**: Change exporter to `awsxray`
- **Google Cloud Trace**: Use `cloudtrace` exporter
- **DataDog/New Relic**: Use vendor-specific exporters

### 4. Environment Variables

```env
# Production
SERVER_APP_OTEL_ENABLED=true
SERVER_APP_OTEL_SERVICE_NAME=go_app_base_prod
SERVER_APP_JAEGER_ENDPOINT=otel-collector:4318  # Use collector sidecar
```

### 5. Security

- **Don't log sensitive data** in span attributes (passwords, tokens, PII)
- Use **sampling** to reduce costs
- Enable **TLS** for collector communication in production

---

## Next Steps

1. **Add metrics**: OpenTelemetry also supports metrics (counters, histograms)
2. **Add logs**: Correlate logs with traces using trace IDs
3. **Custom dashboards**: Create Grafana dashboards for metrics
4. **Alerting**: Set up alerts for high error rates or latency

---

## References

- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)
- [OpenTelemetry Best Practices](https://opentelemetry.io/docs/concepts/instrumentation/manual/)
