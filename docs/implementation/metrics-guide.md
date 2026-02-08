# OpenTelemetry Metrics Guide

This guide explains how to use OpenTelemetry metrics in the Go App Base application.

## Architecture Overview

The metrics implementation follows these key principles:
- **Non-blocking I/O**: All metrics use async aggregation and batch export
- **Vendor-agnostic**: Compatible with any OTLP-compliant backend (Grafana, NewRelic, DataDog, Dynatrace, Jaeger)
- **Zero performance impact**: Background processing with configurable intervals

## Automatic HTTP Metrics

When `SERVER_APP_OTEL_ENABLED=true`, the following HTTP metrics are automatically collected:

### Metric Types

1. **http.server.request.count** (Counter)
   - Description: Total number of HTTP requests
   - Unit: `{request}`
   - Attributes: `http.method`, `http.route`, `http.status_code`

2. **http.server.request.duration** (Histogram)
   - Description: HTTP request duration
   - Unit: `ms` (milliseconds)
   - Attributes: `http.method`, `http.route`, `http.status_code`

3. **http.server.active_requests** (UpDownCounter)
   - Description: Number of active HTTP requests
   - Unit: `{request}`
   - Attributes: `http.method`, `http.route`

4. **http.server.request.size** (Histogram)
   - Description: HTTP request body size
   - Unit: `By` (bytes)
   - Attributes: `http.method`, `http.route`

5. **http.server.response.size** (Histogram)
   - Description: HTTP response body size
   - Unit: `By` (bytes)
   - Attributes: `http.method`, `http.route`, `http.status_code`

## Custom Application Metrics

### Creating Custom Metrics

Use the `CustomMetrics` helper to create application-specific metrics:

```go
package usecases

import (
    "context"
    "github.com/refortunato/go_app_base/internal/shared/observability"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/metric"
)

type CreateExampleUseCase struct {
    repository repositories.ExampleRepository
    metrics    *observability.CustomMetrics
}

func NewCreateExampleUseCase(repo repositories.ExampleRepository) *CreateExampleUseCase {
    return &CreateExampleUseCase{
        repository: repo,
        metrics:    observability.NewCustomMetrics("example_module"),
    }
}
```

### Counter (Monotonic Increment)

Use for counting events that only increase (processed items, errors, etc.):

```go
func (uc *CreateExampleUseCase) Execute(ctx context.Context, input InputDTO) (OutputDTO, error) {
    // Create counter
    processedCounter, err := uc.metrics.Counter(
        "examples.processed.count",
        "Total number of examples processed",
        "{example}",
    )
    if err != nil {
        return OutputDTO{}, err
    }

    // Business logic...
    entity, err := uc.repository.Save(ctx, example)
    if err != nil {
        return OutputDTO{}, err
    }

    // Increment counter (non-blocking)
    processedCounter.Add(ctx, 1,
        metric.WithAttributes(
            attribute.String("status", "success"),
        ),
    )

    return OutputDTO{}, nil
}
```

### UpDownCounter (Can Increment and Decrement)

Use for values that can go up and down (active connections, queue size, etc.):

```go
func (svc *ConnectionService) OpenConnection(ctx context.Context) error {
    activeConns, _ := svc.metrics.UpDownCounter(
        "connections.active",
        "Number of active connections",
        "{connection}",
    )

    // Increment on open
    activeConns.Add(ctx, 1)
    
    // Your logic...
    
    return nil
}

func (svc *ConnectionService) CloseConnection(ctx context.Context) error {
    activeConns, _ := svc.metrics.UpDownCounter(
        "connections.active",
        "Number of active connections",
        "{connection}",
    )

    // Decrement on close
    activeConns.Add(ctx, -1)
    
    return nil
}
```

### Histogram (Distribution of Values)

Use for measuring distributions (response times, payload sizes, etc.):

```go
func (uc *ProcessDataUseCase) Execute(ctx context.Context, data []byte) error {
    // Create histogram
    processingTime, _ := uc.metrics.Histogram(
        "data.processing.duration",
        "Time taken to process data",
        "ms",
    )

    start := time.Now()
    
    // Business logic...
    err := uc.processor.Process(data)
    
    duration := float64(time.Since(start).Milliseconds())
    
    // Record duration (non-blocking)
    processingTime.Record(ctx, duration,
        metric.WithAttributes(
            attribute.Int("data_size", len(data)),
        ),
    )

    return err
}
```

### Gauge (Asynchronous Observation)

Use for values that are observed asynchronously (memory usage, CPU, cache size, etc.):

```go
func (svc *CacheService) RegisterMetrics() error {
    // Gauge callback is executed asynchronously by OpenTelemetry
    err := svc.metrics.Gauge(
        "cache.size",
        "Current number of items in cache",
        "{item}",
        func(ctx context.Context, observer metric.Int64Observer) error {
            // This callback is non-blocking and called periodically
            size := svc.cache.Len()
            observer.Observe(int64(size))
            return nil
        },
    )
    return err
}

// For float values (percentages, ratios, etc.)
func (svc *MemoryService) RegisterMetrics() error {
    err := svc.metrics.FloatGauge(
        "memory.usage.percent",
        "Memory usage percentage",
        "%",
        func(ctx context.Context, observer metric.Float64Observer) error {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            percent := float64(m.Alloc) / float64(m.Sys) * 100
            observer.Observe(percent)
            return nil
        },
    )
    return err
}
```

## Configuration

### Environment Variables

```bash
# Enable/disable metrics
SERVER_APP_OTEL_ENABLED=true

# Service name (used in all metrics)
SERVER_APP_OTEL_SERVICE_NAME=go_app_base

# OTLP endpoint (works with any OTLP-compatible backend)
SERVER_APP_JAEGER_ENDPOINT=jaeger:4318

# Metric export interval in seconds (default: 10s)
# Lower values = more frequent exports, higher overhead
# Higher values = less overhead, delayed visibility
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=10

# Export timeout (default: 30s)
SERVER_APP_OTEL_EXPORT_TIMEOUT=30
```

### Performance Tuning

**For high-traffic applications:**
```bash
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=30  # Export every 30 seconds
```

**For low-latency observability:**
```bash
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=5   # Export every 5 seconds
```

## Backend Configuration

### Grafana Cloud / Grafana Agent

```yaml
# docker-compose.yaml or Kubernetes config
environment:
  - SERVER_APP_JAEGER_ENDPOINT=grafana-agent:4318
```

Grafana automatically ingests OTLP metrics.

### New Relic

```bash
SERVER_APP_JAEGER_ENDPOINT=otlp.nr-data.net:4318
```

Add New Relic API key via OTLP headers (requires custom exporter config).

### DataDog

```bash
SERVER_APP_JAEGER_ENDPOINT=datadog-agent:4318
```

DataDog agent must have OTLP receiver enabled.

### Dynatrace

```bash
SERVER_APP_JAEGER_ENDPOINT=your-environment.live.dynatrace.com:4318
```

Requires Dynatrace API token in headers.

### Jaeger (for testing)

```bash
SERVER_APP_JAEGER_ENDPOINT=jaeger:4318
```

Jaeger v1.35+ supports OTLP metrics ingestion.

## Querying Metrics

### PromQL (Grafana/Prometheus)

```promql
# Request rate
rate(http_server_request_count[5m])

# Average request duration
histogram_quantile(0.95, rate(http_server_request_duration_bucket[5m]))

# Active requests
http_server_active_requests

# Error rate
rate(http_server_request_count{http_status_code=~"5.."}[5m])
```

### NRQL (New Relic)

```sql
SELECT rate(sum(http.server.request.count), 1 minute) FROM Metric

SELECT percentile(http.server.request.duration, 95) FROM Metric
```

## Best Practices

1. **Avoid high-cardinality attributes**: Don't use user IDs, UUIDs, or timestamps as attributes
   ```go
   // ❌ Bad - creates too many unique metric series
   counter.Add(ctx, 1, metric.WithAttributes(
       attribute.String("user_id", userID),
   ))
   
   // ✅ Good - use categories
   counter.Add(ctx, 1, metric.WithAttributes(
       attribute.String("user_type", userType),
   ))
   ```

2. **Reuse metric instruments**: Create once, use many times
   ```go
   // ✅ Create in constructor
   type UseCase struct {
       counter metric.Int64Counter
   }
   
   func NewUseCase(metrics *observability.CustomMetrics) *UseCase {
       counter, _ := metrics.Counter("processed.count", "...", "{item}")
       return &UseCase{counter: counter}
   }
   ```

3. **Use appropriate units**: Follow OpenTelemetry semantic conventions
   - Time: `ms` (milliseconds), `s` (seconds)
   - Size: `By` (bytes), `KiBy` (kibibytes)
   - Count: `{item}`, `{request}`, `{error}`
   - Percentage: `%`

4. **Don't block on metrics**: All operations are already non-blocking, but avoid complex logic in gauge callbacks

5. **Monitor your monitoring**: Keep an eye on metric export failures in logs

## Troubleshooting

### Metrics not appearing in backend

1. Check if observability is enabled:
   ```bash
   SERVER_APP_OTEL_ENABLED=true
   ```

2. Verify endpoint is reachable:
   ```bash
   curl -v http://jaeger:4318/v1/metrics
   ```

3. Check application logs for export errors

4. Verify backend is configured to receive OTLP metrics

### High memory usage

Increase export interval to reduce buffering:
```bash
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=30
```

### Delayed metric visibility

Decrease export interval:
```bash
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=5
```

## Examples

See example module implementations:
- [health/core/application/usecases/check_health.go](../../internal/health/core/application/usecases/check_health.go)
- [example/core/application/usecases/create_example.go](../../internal/example/core/application/usecases/create_example.go)

## Related Documentation

- [Observability Guide](./observability-guide.md) - Complete observability setup
- [OpenTelemetry Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/)
