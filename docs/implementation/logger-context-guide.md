# Logger with Context: Implementation Guide

## Overview

All logger methods now require `context.Context` as the first parameter to automatically extract and propagate trace information from OpenTelemetry.

This enhancement provides:
- **Automatic trace correlation** across all log entries
- **Distributed request tracking** through the entire application stack
- **Enhanced observability** with zero manual effort
- **Preparation for future context values** (userId, tenantId, requestId, etc.)

---

## Core Changes

### Logger Interface

All logging methods now accept `context.Context` as the first parameter:

```go
type Logger interface {
    Debug(ctx context.Context, message string, customFields ...CustomFields)
    Info(ctx context.Context, message string, customFields ...CustomFields)
    Warn(ctx context.Context, message string, customFields ...CustomFields)
    Error(ctx context.Context, message string, customFields ...CustomFields)
    With(fields CustomFields) Logger  // Unchanged - creates logger with persistent fields
}
```

### Automatic Context Extraction

The logger automatically extracts from `context.Context`:
- **`traceId`**: OpenTelemetry trace ID (correlates logs across distributed systems)
- **`spanId`**: OpenTelemetry span ID (identifies specific operation within trace)
- **Future**: Custom values like `userId`, `requestId`, `tenantId`

---

## Usage Patterns by Layer

### 1. Controllers (HTTP Layer)

Extract context from `WebContext` and pass it to logger:

```go
package controllers

func (c *ExampleController) GetExample(webCtx webcontext.WebContext) {
    // Extract context from HTTP request
    ctx := webCtx.GetContext()
    
    // All logger calls use this context
    logger.Info(ctx, "Processing request", logger.CustomFields{
        "exampleId": webCtx.Param("id"),
        "endpoint":  "GET /examples/:id",
    })
    
    // Pass context to use cases
    output, err := c.useCase.Execute(ctx, input)
    if err != nil {
        logger.Error(ctx, "Failed to execute use case", logger.CustomFields{
            "error": err.Error(),
        })
        return
    }
    
    logger.Info(ctx, "Request completed successfully")
    webCtx.JSON(http.StatusOK, output)
}
```

**Key points:**
- Extract `ctx` once at the beginning
- Use same `ctx` for all logger calls in the handler
- Pass `ctx` to use cases for trace propagation

---

### 2. Use Cases (Application Layer)

Use cases already receive context as a parameter:

```go
package usecases

func (uc *GetExampleUseCase) Execute(ctx context.Context, input InputDTO) (OutputDTO, error) {
    // Context already available - just use it
    logger.Info(ctx, "Executing GetExample use case", logger.CustomFields{
        "exampleId": input.Id,
    })
    
    // Pass context to repositories
    entity, err := uc.repository.FindById(ctx, input.Id)
    if err != nil {
        logger.Error(ctx, "Repository error", logger.CustomFields{
            "exampleId": input.Id,
            "error":     err.Error(),
        })
        return nil, err
    }
    
    logger.Debug(ctx, "Entity retrieved from repository")
    return mapToDTO(entity), nil
}
```

**Key points:**
- Use `ctx` parameter directly
- All logs within the same request share the same traceId
- Pass `ctx` to repository calls

---

### 3. Repositories (Infrastructure Layer)

If repositories log (optional), they should accept context:

```go
package repositories

func (r *ExampleMySQLRepository) FindById(ctx context.Context, id string) (*entities.Example, error) {
    logger.Debug(ctx, "Querying database", logger.CustomFields{
        "exampleId": id,
        "table":     "examples",
    })
    
    // Use observability helpers that already accept context
    row := observability.TraceQueryRow(ctx, r.db, "SELECT ...", id)
    
    // Handle result...
    return entity, nil
}
```

**Note:** Most repositories don't need explicit logging since `observability.TraceQuery` already creates spans.

---

### 4. Bootstrap / Initialization (No HTTP Context)

For logs during application startup, use `context.Background()`:

```go
package main

func main() {
    cfg := configs.LoadConfig(".")
    
    // Use Background context for initialization logs
    ctx := context.Background()
    
    logger.Info(ctx, "Application starting", logger.CustomFields{
        "version": cfg.ImageVersion,
        "env":     cfg.Environment,
    })
    
    db := setupDatabase(cfg)
    logger.Info(ctx, "Database connected successfully")
    
    // ... rest of initialization
}
```

**When to use `context.Background()`:**
- Application startup
- Background workers (without HTTP request)
- Scheduled jobs
- Standalone scripts

---

## Trace Correlation Examples

### Example 1: HTTP Request Flow

**Request:** `GET /examples/123`

**Logs generated:**

```json
{
  "timestamp": "2026-02-08T10:15:30.123456Z",
  "level": "INFO",
  "msg": "Processing GetExample request",
  "imageName": "go_app_base",
  "imageVersion": "1.0.0",
  "custom": {
    "traceId": "4bf92f3577b34da6a3ce929d0e0e4736",
    "spanId": "00f067aa0ba902b7",
    "exampleId": "123",
    "endpoint": "GET /examples/:id"
  }
}
```

```json
{
  "timestamp": "2026-02-08T10:15:30.234567Z",
  "level": "INFO",
  "msg": "Executing GetExample use case",
  "imageName": "go_app_base",
  "imageVersion": "1.0.0",
  "custom": {
    "traceId": "4bf92f3577b34da6a3ce929d0e0e4736",  // Same trace!
    "spanId": "349f8a7b2d1c5e92",  // Different span
    "exampleId": "123"
  }
}
```

```json
{
  "timestamp": "2026-02-08T10:15:30.345678Z",
  "level": "INFO",
  "msg": "Request completed successfully",
  "imageName": "go_app_base",
  "imageVersion": "1.0.0",
  "custom": {
    "traceId": "4bf92f3577b34da6a3ce929d0e0e4736",  // Same trace!
    "spanId": "00f067aa0ba902b7"  // Back to controller span
  }
}
```

**Analysis:**
- All logs share the **same `traceId`** → they belong to the same request
- Each layer has a **different `spanId`** → you can see which layer logged
- You can query logs by `traceId` to see the complete request journey

---

### Example 2: Error Scenario

**Request:** `GET /examples/invalid-id`

```json
{
  "timestamp": "2026-02-08T10:20:15.123456Z",
  "level": "ERROR",
  "msg": "Repository error",
  "imageName": "go_app_base",
  "imageVersion": "1.0.0",
  "custom": {
    "traceId": "7d2a9f5e8c1b4a3d6f0e2b9c5a8d1f4e",
    "spanId": "2c5d8f1a3b9e7d4c",
    "exampleId": "invalid-id",
    "error": "sql: no rows in result set"
  }
}
```

**Query logs by traceId** to see:
1. Initial request log (controller)
2. Use case execution log
3. Error log (with full context)
4. Error response log

---

## Integration with Jaeger

### Viewing Traces in Jaeger UI

1. **Access Jaeger**: `http://localhost:16686`

2. **Search by Service**: Select `go_app_base`

3. **Find Trace**: Click on a trace to see the timeline

4. **Correlate with Logs**:
   - Copy `traceId` from Jaeger
   - Query logs: `grep "traceId": "4bf92f..." logs.json`
   - See all log entries for that request

### Example Trace View

```
GET /examples/123
├── Span: GET /examples/:id (Controller)
│   ├── Log: "Processing GetExample request"
│   └── Span: GetExampleUseCase.Execute
│       ├── Log: "Executing GetExample use case"
│       └── Span: SELECT * FROM examples WHERE id = ?
│           └── Log: "Querying database"
```

---

## Migration Checklist

If you're updating existing code:

- [ ] All controller methods extract `ctx := webCtx.GetContext()` at the start
- [ ] All logger calls in controllers use `logger.Info(ctx, ...)`
- [ ] All use cases receive `context.Context` as first parameter
- [ ] All logger calls in use cases use the received `ctx`
- [ ] Repositories accept `context.Context` if they log
- [ ] Initialization logs use `context.Background()`
- [ ] No logger calls without context remain in codebase

---

## Best Practices

### ✅ DO

```go
// Extract context once
ctx := webCtx.GetContext()

// Use it consistently
logger.Info(ctx, "Starting operation")
result, err := service.Execute(ctx, input)
logger.Info(ctx, "Operation completed")
```

### ✅ DO

```go
// Use Background for non-request logs
ctx := context.Background()
logger.Info(ctx, "Application starting")
```

### ✅ DO

```go
// Pass context through the call stack
func (uc *UseCase) Execute(ctx context.Context, input InputDTO) (OutputDTO, error) {
    logger.Info(ctx, "Use case started")
    entity, err := uc.repo.FindById(ctx, input.Id)
    // ...
}
```

### ❌ DON'T

```go
// Don't call logger without context
logger.Info("This won't compile")  // ❌ Missing context parameter
```

### ❌ DON'T

```go
// Don't create new context unnecessarily
ctx := context.Background()  // ❌ Use received context instead
logger.Info(ctx, "...")
```

### ❌ DON'T

```go
// Don't pass nil context
logger.Info(nil, "This will work but lose trace info")  // ❌ Prefer context.Background()
```

---

## Troubleshooting

### Issue: Logs don't have traceId

**Possible causes:**
1. Not using `ctx` from HTTP request → Use `webCtx.GetContext()`
2. Passing `nil` or `context.Background()` in HTTP handlers
3. OpenTelemetry not initialized → Check main.go initialization

**Solution:**
```go
// In controllers
ctx := webCtx.GetContext()  // ✅ Extracts trace from HTTP request
logger.Info(ctx, "...")

// NOT this:
ctx := context.Background()  // ❌ Creates new context without trace
```

---

### Issue: Different traceIds for same request

**Cause:** Creating new context instead of propagating existing one

**Solution:**
```go
// ✅ Correct - propagate context
func (uc *UseCase) Execute(ctx context.Context, input InputDTO) {
    logger.Info(ctx, "...")  // Uses received context
}

// ❌ Wrong - creates new context
func (uc *UseCase) Execute(ctx context.Context, input InputDTO) {
    newCtx := context.Background()  // Don't do this!
    logger.Info(newCtx, "...")
}
```

---

## Future Enhancements

The context infrastructure is prepared for:

### 1. User ID Tracking

```go
// In authentication middleware
ctx = context.WithValue(ctx, "userId", user.ID)

// Logger will automatically extract
logger.Info(ctx, "User action")
// Output: { ..., "custom": { "traceId": "...", "userId": "user-123" } }
```

### 2. Request ID

```go
// In middleware
requestId := uuid.New()
ctx = context.WithValue(ctx, "requestId", requestId)

// All logs include requestId
```

### 3. Tenant ID (Multi-tenancy)

```go
ctx = context.WithValue(ctx, "tenantId", tenant.ID)
logger.Info(ctx, "Tenant operation")
```

**To enable:** Update `ExtractCustomContextFields()` in [internal/shared/logger/context_utils.go](internal/shared/logger/context_utils.go)

---

## Summary

- ✅ All logger methods require `context.Context` as first parameter
- ✅ TraceId/SpanId extracted automatically from OpenTelemetry
- ✅ Controllers extract context from HTTP request
- ✅ Use cases receive and propagate context
- ✅ Use `context.Background()` for non-request logs
- ✅ Logs are automatically correlated in Jaeger and log aggregators
- ✅ Ready for future context values (userId, requestId, etc.)

**Result:** Complete observability with zero manual effort! 🎉
