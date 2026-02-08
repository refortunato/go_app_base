# OpenTelemetry Collector Architecture

## Overview

This project uses **OpenTelemetry Collector** as a central observability hub that receives traces and metrics from the application and forwards them to specialized backends (Jaeger for traces, Prometheus for metrics).

---

## Architecture Diagram

```
┌────────────────────────────────────────────────────────────────────┐
│                         GO APPLICATION                             │
│                         (Port 8080)                                │
│                                                                    │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  OpenTelemetry SDK                                          │   │
│  │  ┌──────────────────┐      ┌───────────────────┐            │   │
│  │  │  TracerProvider  │      │  MeterProvider    │            │   │
│  │  │  (Traces)        │      │  (Metrics)        │            │   │
│  │  └────────┬─────────┘      └─────────┬─────────┘            │   │
│  │           │                          │                      │   │
│  │           └───────────┬──────────────┘                      │   │
│  │                       │                                     │   │
│  │               ┌───────▼────────┐                            │   │
│  │               │  OTLP Exporter │                            │   │
│  │               │  (HTTP)        │                            │   │
│  │               └───────┬────────┘                            │   │
│  └───────────────────────┼─────────────────────────────────────┘   │
└──────────────────────────┼─────────────────────────────────────────┘
                           │
                           │ OTLP/HTTP
                           │ (Traces + Metrics)
                           │ Port 4318
                           │
                           ▼
┌────────────────────────────────────────────────────────────────────┐
│              OPENTELEMETRY COLLECTOR                               │
│              (otel-collector:4318)                                 │
│                                                                    │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  RECEIVERS                                                   │  │
│  │  ┌────────────────┐                                          │  │
│  │  │  OTLP/HTTP     │  Receives traces & metrics from app      │  │
│  │  │  Port: 4318    │                                          │  │
│  │  └────────┬───────┘                                          │  │
│  └───────────┼──────────────────────────────────────────────────┘  │
│              │                                                     │
│  ┌───────────▼──────────────────────────────────────────────────┐  │
│  │  PROCESSORS                                                  │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐     │  │
│  │  │ Memory       │  │ Resource     │  │ Batch           │     │  │
│  │  │ Limiter      │→ │ Processor    │→ │ Processor       │     │  │
│  │  │ (512 MiB)    │  │ (Add labels) │  │ (10s / 1024)    │     │  │
│  │  └──────────────┘  └──────────────┘  └─────────┬───────┘     │  │
│  └────────────────────────────────────────────────┼─────────────┘  │
│                                                   │                │
│  ┌────────────────────────────────────────────────▼──────────────┐ │
│  │  EXPORTERS                                                    │ │
│  │                                                               │ │
│  │  ┌─────────────────┐              ┌──────────────────────┐    │ │
│  │  │  OTLP/Jaeger    │              │  Prometheus          │    │ │
│  │  │  (Traces only)  │              │  (Metrics only)      │    │ │
│  │  │  Port: 4317     │              │  Port: 8889          │    │ │
│  │  └────────┬────────┘              └──────────┬───────────┘    │ │
│  └───────────┼──────────────────────────────────┼────────────────┘ │
└──────────────┼──────────────────────────────────┼──────────────────┘
               │                                  │
               │ gRPC                             │ HTTP
               │ Traces                           │ Metrics (/metrics)
               │                                  │ Scrape Interval: 10s
               │                                  │
               ▼                                  ▼
┌──────────────────────────┐      ┌────────────────────────────────┐
│       JAEGER             │      │       PROMETHEUS               │
│   (Port 16686 - UI)      │      │   (Port 9090 - UI)             │
│                          │      │                                │
│  ┌────────────────────┐  │      │  ┌───────────────────────────┐ │
│  │  Trace Storage     │  │      │  │  Metrics Storage          │ │
│  │  (In-Memory)       │  │      │  │  (Time Series DB)         │ │
│  └────────────────────┘  │      │  └───────────┬───────────────┘ │
│                          │      │              │                 │
│  ┌────────────────────┐  │      │  ┌───────────▼──────────────┐  │
│  │  Trace Query API   │  │      │  │  PromQL Query Engine     │  │
│  └────────────────────┘  │      │  └──────────────────────────┘  │
└──────────────────────────┘      └───────────────┬────────────────┘
                                                  │
                                                  │ Datasource
                                                  │ http://prometheus:9090
                                                  │
                                                  ▼
                                  ┌────────────────────────────────┐
                                  │       GRAFANA                  │
                                  │   (Port 3000 - UI)             │
                                  │                                │
                                  │  ┌───────────────────────────┐ │
                                  │  │  Dashboards               │ │
                                  │  │  - HTTP Metrics           │ │
                                  │  │  - Request Rate           │ │
                                  │  │  - Error Rate             │ │
                                  │  │  - Latency (P95, P99)     │ │
                                  │  └───────────────────────────┘ │
                                  └────────────────────────────────┘
```

---

## Components

### 1. Go Application (Port 8080)

**Role**: Generates observability data (traces and metrics)

**Libraries**:
- `go.opentelemetry.io/otel` - OpenTelemetry SDK
- `go.opentelemetry.io/otel/exporters/otlp/otlptracehttp` - Trace exporter
- `go.opentelemetry.io/otel/exporters/otlp/otlpmetrichttp` - Metric exporter

**Configuration**:
```env
SERVER_APP_OTEL_ENABLED=true
SERVER_APP_OTEL_SERVICE_NAME=go_app_base
SERVER_APP_JAEGER_ENDPOINT=otel-collector:4318  # Points to Collector
```

**What it exports**:
- **Traces**: HTTP requests, database queries, use case execution
- **Metrics**: Request count, duration, active requests, request/response size

---

### 2. OpenTelemetry Collector (Ports 4318, 8889)

**Role**: Central observability hub - receives, processes, and routes telemetry data

**Image**: `otel/opentelemetry-collector-contrib:0.99.0`

**Configuration**: `otel-collector-config.yaml`

#### Receivers
- **OTLP/HTTP (4318)**: Receives traces and metrics from Go app
- **OTLP/gRPC (4317)**: Alternative protocol (not used currently)

#### Processors
- **Memory Limiter**: Prevents OOM by dropping data when memory > 512 MiB
- **Resource Processor**: Adds custom attributes (service.namespace, deployment.environment)
- **Batch Processor**: Groups data (10s timeout, 1024 batch size) to reduce network overhead

#### Exporters
- **OTLP/Jaeger (4317)**: Forwards traces to Jaeger via gRPC
- **Prometheus (8889)**: Exposes metrics in Prometheus format at `/metrics` endpoint
- **Logging**: Outputs sampled data to logs for debugging

#### Pipelines
```yaml
traces:  OTLP → Memory Limiter → Resource → Batch → Jaeger
metrics: OTLP → Memory Limiter → Resource → Batch → Prometheus Format
```

---

### 3. Jaeger (Port 16686)

**Role**: Distributed tracing backend - stores and visualizes traces

**Image**: `jaegertracing/all-in-one:1.53`

**Receives from**: OpenTelemetry Collector (port 4317 gRPC)

**Features**:
- Trace search and filtering
- Service dependency graph
- Latency analysis
- Error tracking

**Access**: http://localhost:16686

---

### 4. Prometheus (Port 9090)

**Role**: Metrics aggregation and storage (time-series database)

**Image**: `prom/prometheus:v2.55.1`

**Configuration**: `prometheus.yml`

**Scrape Configuration**:
```yaml
- job_name: 'otel-collector'
  static_configs:
    - targets: ['otel-collector:8889']
  scrape_interval: 10s
```

**What it scrapes**:
- Application metrics from Collector (go_app_base.http.server.*)
- OpenTelemetry internal metrics (otel.*)
- Prometheus self-monitoring metrics

**Access**: http://localhost:9090

---

### 5. Grafana (Port 3000)

**Role**: Metrics visualization and dashboards

**Image**: `grafana/grafana:11.4.0`

**Datasource**: Prometheus (http://prometheus:9090)

**Features**:
- Custom dashboards
- Alerting
- Query builder
- Panel templates

**Default Credentials**:
- Username: `admin`
- Password: `admin`

**Access**: http://localhost:3000

---

## Data Flow

### Traces Flow
```
App → TracerProvider → OTLP Exporter → Collector (4318) 
    → Batch Processor → OTLP/Jaeger Exporter → Jaeger (4317)
    → Trace Storage → Jaeger UI (16686)
```

### Metrics Flow
```
App → MeterProvider → OTLP Exporter → Collector (4318)
    → Batch Processor → Prometheus Exporter (8889)
    → Prometheus Scrape (10s interval) → Time Series DB
    → PromQL Engine → Grafana Dashboards (3000)
```

---

## Key Advantages

### ✅ Centralization
- **Single endpoint**: App sends to Collector only
- **Decoupling**: App doesn't know about Jaeger or Prometheus
- **Flexibility**: Easy to add new backends (Datadog, New Relic, etc.)

### ✅ Performance
- **Batching**: Reduces network overhead (groups 1024 data points every 10s)
- **Buffering**: Collector queues data during backend outages
- **Retry logic**: Automatic retry with exponential backoff

### ✅ Processing
- **Filtering**: Drop unnecessary data
- **Sampling**: Keep only N% of traces
- **Enrichment**: Add custom attributes/labels
- **Transformation**: Convert formats

### ✅ Reliability
- **Memory limits**: Prevents OOM crashes
- **Queue management**: Handles traffic spikes
- **Health checks**: Kubernetes-ready liveness probes

---

## Configuration Details

### Collector Resources

```yaml
memory_limiter:
  limit_mib: 512        # Hard limit
  spike_limit_mib: 128  # Allow temporary spikes

batch:
  timeout: 10s          # Send every 10s
  send_batch_size: 1024 # Or when 1024 items accumulated
```

### Retry Configuration

```yaml
retry_on_failure:
  initial_interval: 5s  # First retry after 5s
  max_interval: 30s     # Max wait between retries
  max_elapsed_time: 300s # Give up after 5 minutes
```

---

## Ports Reference

| Component | Port | Protocol | Purpose |
|-----------|------|----------|---------|
| **App** | 8080 | HTTP | Application API |
| **Collector** | 4318 | HTTP | OTLP receiver (from app) |
| **Collector** | 4317 | gRPC | OTLP receiver (alternative) |
| **Collector** | 8889 | HTTP | Prometheus exporter |
| **Collector** | 8888 | HTTP | Internal metrics |
| **Jaeger** | 4317 | gRPC | OTLP receiver (from collector) |
| **Jaeger** | 16686 | HTTP | Jaeger UI |
| **Prometheus** | 9090 | HTTP | Prometheus UI + API |
| **Grafana** | 3000 | HTTP | Grafana UI |

---

## Environment Variables

### Application

```env
# Enable observability
SERVER_APP_OTEL_ENABLED=true

# Service identification
SERVER_APP_OTEL_SERVICE_NAME=go_app_base

# Collector endpoint (NOT Jaeger directly)
SERVER_APP_JAEGER_ENDPOINT=otel-collector:4318

# Metric export interval
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=10
```

### Docker Compose

```yaml
environment:
  SERVER_APP_OTEL_ENABLED: "true"
  SERVER_APP_OTEL_SERVICE_NAME: go_app_base
  SERVER_APP_JAEGER_ENDPOINT: otel-collector:4318
```

---

## Monitoring the Collector

### Health Check
```bash
curl http://localhost:8888/
```

### Internal Metrics
```bash
curl http://localhost:8888/metrics
```

**Key metrics**:
- `otelcol_receiver_accepted_spans` - Traces received
- `otelcol_receiver_accepted_metric_points` - Metrics received
- `otelcol_exporter_sent_spans` - Traces sent to Jaeger
- `otelcol_processor_batch_batch_send_size` - Batch sizes

### Logs
```bash
docker logs go_app_base_otel_collector
```

---

## Testing the Pipeline

### 1. Generate Traffic
```bash
for i in {1..50}; do 
    curl http://localhost:8080/health
done
```

### 2. Check Collector Reception
```bash
# Verify Collector is receiving data
curl -s http://localhost:8888/metrics | grep otelcol_receiver_accepted
```

Expected output:
```
otelcol_receiver_accepted_spans{...} 50
otelcol_receiver_accepted_metric_points{...} 250
```

### 3. Check Collector Export
```bash
# Verify Collector is exporting to Prometheus
curl -s http://localhost:8889/metrics | grep go_app_base
```

Expected output:
```
go_app_base_http_server_request_count{...} 50
go_app_base_http_server_request_duration_bucket{...} 45
```

### 4. Check Prometheus Scrape
```bash
# Verify Prometheus scraped the metrics
curl -s 'http://localhost:9090/api/v1/query?query=go_app_base_http_server_request_count' | jq
```

### 5. Check Jaeger Traces
```bash
# Access Jaeger UI
open http://localhost:16686
# Search for service: go_app_base
```

### 6. Visualize in Grafana
```bash
# Access Grafana
open http://localhost:3000
# Login: admin / admin
# Add Prometheus datasource: http://prometheus:9090
```

---

## Troubleshooting

### Metrics not appearing in Prometheus

**Check 1**: Collector receiving data?
```bash
docker logs go_app_base_otel_collector | grep -i "metric"
```

**Check 2**: Collector exporting to Prometheus format?
```bash
curl -s http://localhost:8889/metrics | head -20
```

**Check 3**: Prometheus scraping Collector?
```bash
curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.job=="otel-collector")'
```

### Traces not appearing in Jaeger

**Check 1**: Collector receiving traces?
```bash
curl -s http://localhost:8888/metrics | grep accepted_spans
```

**Check 2**: Collector exporting to Jaeger?
```bash
docker logs go_app_base_otel_collector | grep -i "jaeger"
```

**Check 3**: Jaeger receiving spans?
```bash
curl -s http://localhost:16686/api/services | jq
```

### Collector OOM or high memory

**Solution 1**: Reduce batch size
```yaml
batch:
  send_batch_size: 512  # Down from 1024
```

**Solution 2**: Increase memory limit
```yaml
memory_limiter:
  limit_mib: 1024  # Up from 512
```

---

## Scaling Considerations

### Horizontal Scaling
- Run multiple Collector instances behind load balancer
- Each instance processes subset of data
- Shared-nothing architecture

### Vertical Scaling
- Increase memory limit
- Increase batch size
- Add more processors

### Production Recommendations
- Use persistent storage for Jaeger (Cassandra, Elasticsearch)
- Enable sampling in Collector (not all traces)
- Set up alerting on Collector health metrics
- Use Collector in agent mode (sidecar) + gateway mode (centralized)

---

## References

- [OpenTelemetry Collector Documentation](https://opentelemetry.io/docs/collector/)
- [Collector Configuration Reference](https://opentelemetry.io/docs/collector/configuration/)
- [OTLP Protocol Specification](https://opentelemetry.io/docs/specs/otlp/)
- [Prometheus Exporter Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/prometheusexporter)

---

## Next Steps

1. ✅ **Start environment**: `make dev`
2. ✅ **Generate traffic**: Use curl or application
3. ✅ **View traces**: http://localhost:16686
4. ✅ **Query metrics**: http://localhost:9090
5. 📊 **Create dashboards**: http://localhost:3000
6. 🚀 **Add custom metrics**: Use `CustomMetrics` interface
7. ⚙️ **Fine-tune Collector**: Adjust batch sizes, sampling, etc.
