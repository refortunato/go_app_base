# Dynamic Metrics Naming

## Overview

Metrics are automatically prefixed with your application name (configured via `SERVER_APP_NAME`), making it easier to identify metrics in multi-service environments like Prometheus/Grafana.

## Configuration

### Environment Variable

```env
# Application name - used as metric prefix
SERVER_APP_NAME=ms-registration
```

**Result**: All HTTP metrics will be prefixed with `ms_registration.*`

### Naming Conventions

The system automatically normalizes application names to be Prometheus-compliant:

| App Name | Metric Prefix | Example Metric |
|----------|---------------|----------------|
| `go_app_base` | `go_app_base` | `go_app_base.http.server.request.count` |
| `ms-registration` | `ms_registration` | `ms_registration.http.server.request.count` |
| `user-service` | `user_service` | `user_service.http.server.request.count` |
| `order_api` | `order_api` | `order_api.http.server.request.count` |

**Rules**:
- Hyphens (`-`) are converted to underscores (`_`)
- Lowercase is preserved
- No other transformations applied

## Available HTTP Metrics

All HTTP metrics follow the pattern: `<app_name>.http.server.*`

### 1. Request Count (Counter)
```promql
<app_name>.http.server.request.count
```

**Labels**:
- `http.method` - HTTP method (GET, POST, etc.)
- `http.route` - Endpoint route (/users/:id)
- `http.status_code` - HTTP status code (200, 404, 500)

**Example**:
```promql
# Total requests to /users endpoint
ms_registration.http.server.request.count{http_route="/users"}

# Failed requests (5xx errors)
ms_registration.http.server.request.count{http_status_code=~"5.."}
```

### 2. Request Duration (Histogram)
```promql
<app_name>.http.server.request.duration
```

**Unit**: milliseconds (ms)

**Labels**: Same as request count

**Example**:
```promql
# P95 latency for /orders endpoint
histogram_quantile(0.95, 
  sum by(le, http_route) (
    rate(ms_registration.http.server.request.duration_bucket{http_route="/orders"}[5m])
  )
)
```

### 3. Active Requests (UpDownCounter)
```promql
<app_name>.http.server.active_requests
```

**Labels**:
- `http.method`
- `http.route`

**Example**:
```promql
# Current active requests per endpoint
ms_registration.http.server.active_requests
```

### 4. Request Size (Histogram)
```promql
<app_name>.http.server.request.size
```

**Unit**: bytes (By)

**Example**:
```promql
# Average request size
rate(ms_registration.http.server.request.size_sum[5m]) 
/ 
rate(ms_registration.http.server.request.size_count[5m])
```

### 5. Response Size (Histogram)
```promql
<app_name>.http.server.response.size
```

**Unit**: bytes (By)

**Example**:
```promql
# P99 response size
histogram_quantile(0.99, 
  rate(ms_registration.http.server.response.size_bucket[5m])
)
```

## How It Works

### Architecture

```
Environment Variable      Normalization        Metric Creation
     ↓                         ↓                     ↓
SERVER_APP_NAME      normalizeMetricPrefix()   meter.Int64Counter()
   (ms-registration)      (ms_registration)       (ms_registration.http.server.request.count)
```

### Code Flow

1. **Configuration Load** ([configs/config.go](../../configs/config.go)):
   ```go
   AppName: getEnv("SERVER_APP_NAME", "go_app_base")
   ```

2. **Server Creation** ([cmd/server/main.go](../../cmd/server/main.go)):
   ```go
   server.NewGinServerWithRoutes(
       cfg.WebServerPort,
       infraWeb.RegisterRoutes(c),
       cfg.OtelServiceName,
       cfg.AppName,  // 👈 Passed to factory
       cfg.OtelEnabled,
   )
   ```

3. **Middleware Registration** ([factory.go](../../internal/shared/web/server/factory.go)):
   ```go
   router.Use(observability.MetricsMiddleware(serviceName, appName))
   ```

4. **Metric Naming** ([metrics_middleware.go](../../internal/shared/observability/metrics_middleware.go)):
   ```go
   func MetricsMiddleware(serviceName, appName string) gin.HandlerFunc {
       metricPrefix := normalizeMetricPrefix(appName)
       requestCounter, _ := meter.Int64Counter(
           metricPrefix + ".http.server.request.count",
           // ...
       )
   }
   ```

## Changing Your Application Name

### Step 1: Update Environment Variable

```bash
# cmd/server/.env
SERVER_APP_NAME=ms-registration
```

### Step 2: Rebuild & Restart

```bash
make down
make dev
```

### Step 3: Verify in Prometheus

```bash
# Access Prometheus UI
open http://localhost:9090

# Query your new metrics
ms_registration.http.server.request.count
```

## Multi-Service Monitoring

When running multiple services, each will have its own metric prefix:

```promql
# Service 1 (ms-registration)
ms_registration.http.server.request.count

# Service 2 (user-service)  
user_service.http.server.request.count

# Service 3 (order-api)
order_api.http.server.request.count
```

### Aggregation Across Services

```promql
# Total requests across all services
sum({__name__=~".*http.server.request.count"})

# Per-service comparison
sum by(__name__) ({__name__=~".*http.server.request.count"})
```

## Custom Metrics

When creating custom metrics, follow the same naming pattern:

```go
import "github.com/refortunato/go_app_base/internal/shared/observability"

func (uc *MyUseCase) Execute() {
    customMetrics := observability.NewCustomMetrics("my-service")
    
    // Use app name as prefix manually
    counter, _ := customMetrics.Counter(
        "ms_registration.business.orders.created",  // 👈 Use app prefix
        "Number of orders created",
        "{order}",
    )
    
    counter.Add(ctx, 1)
}
```

**Better approach** - Create a helper:

```go
func (cm *CustomMetrics) CounterWithPrefix(appName, name, description, unit string) {
    prefix := normalizeMetricPrefix(appName)
    return cm.Counter(prefix + "." + name, description, unit)
}
```

## Grafana Dashboard

### Template Variable

Create a variable to switch between services:

**Name**: `service`  
**Type**: Query  
**Query**: 
```promql
label_values({__name__=~".*http.server.request.count"}, __name__)
```

### Panel Queries

Use the variable in queries:

```promql
# Request rate
rate($service{http_route="$endpoint"}[5m])

# Error rate
sum(rate($service{http_status_code=~"5.."}[5m])) 
/ 
sum(rate($service[5m]))
```

## Best Practices

### ✅ Do

- Use descriptive, meaningful application names
- Keep names short and consistent
- Use lowercase with hyphens or underscores
- Document your naming convention in README

### ❌ Don't

- Use spaces in application names
- Use special characters (!, @, #, etc.)
- Change names frequently in production
- Use different names across environments (dev/staging/prod)

## Testing

### Verify Metric Prefix

```bash
# Generate traffic
for i in {1..50}; do 
    curl http://localhost:8080/health
done

# Wait for export (10s default)
sleep 15

# Check Jaeger metrics endpoint
curl -s http://localhost:14269/metrics | grep "$(echo $SERVER_APP_NAME | tr '-' '_')"

# Should show metrics like:
# ms_registration_http_server_request_count{...} 50
```

### Prometheus Query Test

```bash
# Query via API
curl -s 'http://localhost:9090/api/v1/query?query=ms_registration_http_server_request_count' | jq

# Expected result:
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "http_method": "GET",
          "http_route": "/health",
          "http_status_code": "200"
        },
        "value": [1707321600, "50"]
      }
    ]
  }
}
```

## Troubleshooting

### Metrics still using old prefix

**Cause**: Application not restarted or old metrics cached

**Solution**:
```bash
make down
docker system prune -f
make dev
```

### Metrics with both old and new prefixes

**Cause**: Multiple application instances with different names

**Solution**: Ensure all instances use the same `SERVER_APP_NAME`

### Invalid metric names in Prometheus

**Cause**: Special characters in app name

**Solution**: Use only alphanumeric, hyphens, and underscores

---

## Summary

- ✅ Metrics automatically prefixed with `SERVER_APP_NAME`
- ✅ Hyphens converted to underscores (Prometheus-compliant)
- ✅ Easy to identify metrics in multi-service environments
- ✅ No code changes needed when renaming application
- ✅ Consistent across all HTTP metrics

**Next**: [Metrics Troubleshooting Guide](../METRICS_TROUBLESHOOTING.md)
