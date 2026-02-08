# Métricas OpenTelemetry - Exemplos de Uso

## Visão Geral

Este documento contém exemplos práticos de uso de métricas OpenTelemetry no go_app_base.

## Métricas Automáticas HTTP

Quando `SERVER_APP_OTEL_ENABLED=true`, as seguintes métricas são coletadas automaticamente para todas as requisições HTTP:

- `http.server.request.count` - Total de requisições
- `http.server.request.duration` - Duração das requisições (histogram)
- `http.server.active_requests` - Requisições ativas no momento
- `http.server.request.size` - Tamanho do body da requisição
- `http.server.response.size` - Tamanho do body da resposta

## Métricas Customizadas - Exemplo Real

### 1. Health Check com Counter

Arquivo: `internal/health/core/application/usecases/health_check.go`

```go
type HealthCheckUseCase struct {
    healthRepository repositories.HealthRepository
    metrics          *observability.CustomMetrics
    healthCounter    metric.Int64Counter
}

func NewHealthCheckUseCase(healthRepository repositories.HealthRepository) *HealthCheckUseCase {
    metrics := observability.NewCustomMetrics("health_module")
    
    // Criar counter uma vez (reutilizar em todas as chamadas)
    healthCounter, _ := metrics.Counter(
        "health.check.count",
        "Total number of health checks performed",
        "{check}",
    )
    
    return &HealthCheckUseCase{
        healthRepository: healthRepository,
        metrics:          metrics,
        healthCounter:    healthCounter,
    }
}

func (u *HealthCheckUseCase) Execute() (*HealthCheckOutputDTO, error) {
    ctx := context.Background()
    
    err := u.healthRepository.CheckDatabaseConnection()
    
    // Registrar métrica (não bloqueante)
    status := "success"
    if err != nil {
        status = "failure"
    }
    
    u.healthCounter.Add(ctx, 1,
        metric.WithAttributes(
            attribute.String("status", status),
        ),
    )
    
    if err != nil {
        return nil, err
    }

    return &HealthCheckOutputDTO{Status: "OK"}, nil
}
```

### 2. Use Case com Múltiplas Métricas

Arquivo: `internal/example/core/application/usecases/metrics_demo.go`

```go
type CreateExampleMetricsDemo struct {
    repository repositories.ExampleRepository
    
    // Instrumentos de métricas (criar uma vez, reutilizar sempre)
    metrics           *observability.CustomMetrics
    creationCounter   metric.Int64Counter      // Total criado
    creationDuration  metric.Float64Histogram  // Tempo de criação
    activeCreations   metric.Int64UpDownCounter // Operações ativas
}

func NewCreateExampleMetricsDemo(repo repositories.ExampleRepository) *CreateExampleMetricsDemo {
    metrics := observability.NewCustomMetrics("example_module")
    
    creationCounter, _ := metrics.Counter(
        "examples.created.total",
        "Total number of examples created",
        "{example}",
    )
    
    creationDuration, _ := metrics.Histogram(
        "examples.creation.duration",
        "Time taken to create an example",
        "ms",
    )
    
    activeCreations, _ := metrics.UpDownCounter(
        "examples.creation.active",
        "Number of in-progress example creations",
        "{operation}",
    )
    
    return &CreateExampleMetricsDemo{
        repository:        repo,
        metrics:           metrics,
        creationCounter:   creationCounter,
        creationDuration:  creationDuration,
        activeCreations:   activeCreations,
    }
}

func (uc *CreateExampleMetricsDemo) Execute(ctx context.Context, name string) (*entities.Example, error) {
    // Incrementar operações ativas
    uc.activeCreations.Add(ctx, 1)
    defer func() {
        // Decrementar ao finalizar
        uc.activeCreations.Add(ctx, -1)
    }()
    
    start := time.Now()
    
    // Lógica de negócio...
    example, err := entities.NewExample(name)
    if err != nil {
        uc.creationCounter.Add(ctx, 1,
            metric.WithAttributes(attribute.String("status", "validation_error")),
        )
        return nil, err
    }
    
    savedExample, err := uc.repository.Save(ctx, example)
    if err != nil {
        uc.creationCounter.Add(ctx, 1,
            metric.WithAttributes(attribute.String("status", "repository_error")),
        )
        return nil, err
    }
    
    // Registrar métricas de sucesso
    duration := float64(time.Since(start).Milliseconds())
    
    uc.creationCounter.Add(ctx, 1,
        metric.WithAttributes(attribute.String("status", "success")),
    )
    
    uc.creationDuration.Record(ctx, duration,
        metric.WithAttributes(attribute.String("operation", "create")),
    )
    
    return savedExample, nil
}
```

### 3. Gauge Assíncrono para Monitoramento

```go
// Callback executado periodicamente de forma assíncrona
func (uc *CreateExampleMetricsDemo) RegisterGaugeMetrics() error {
    return uc.metrics.Gauge(
        "examples.repository.size",
        "Approximate number of examples in repository",
        "{example}",
        func(ctx context.Context, observer metric.Int64Observer) error {
            // IMPORTANTE: Usar valor em cache, não queries pesadas
            count := getCachedCount() // NÃO fazer COUNT(*) aqui
            observer.Observe(count)
            return nil
        },
    )
}
```

## Configuração Recomendada por Ambiente

### Desenvolvimento
```bash
SERVER_APP_OTEL_ENABLED=true
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=5  # 5 segundos para feedback rápido
```

### Staging
```bash
SERVER_APP_OTEL_ENABLED=true
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=10  # 10 segundos (padrão)
```

### Produção - Alta Carga
```bash
SERVER_APP_OTEL_ENABLED=true
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=30  # 30 segundos para menor overhead
```

## Queries de Exemplo (PromQL)

### Taxa de Health Checks com Sucesso
```promql
rate(health_check_count{status="success"}[5m])
```

### Taxa de Falhas
```promql
rate(health_check_count{status="failure"}[5m])
```

### Percentil 95 de Duração HTTP
```promql
histogram_quantile(0.95, rate(http_server_request_duration_bucket[5m]))
```

### Requisições Ativas
```promql
http_server_active_requests
```

### Taxa de Criação de Exemplos
```promql
rate(examples_created_total{status="success"}[5m])
```

### Tempo Médio de Criação
```promql
avg(examples_creation_duration)
```

## Boas Práticas Demonstradas

1. ✅ **Criar instrumentos no construtor** (reutilizar em todas as chamadas)
2. ✅ **Usar atributos para categorizar** (status, operation, etc.)
3. ✅ **Registrar tanto sucesso quanto falha** (visibilidade completa)
4. ✅ **UpDownCounter para operações ativas** (tracking em tempo real)
5. ✅ **Histogram para distribuições** (latências, tamanhos, etc.)
6. ✅ **Gauge para valores observados** (cache size, connections, etc.)
7. ✅ **Todas as operações são não-bloqueantes** (zero impacto em I/O)

## Testando

```bash
# 1. Subir aplicação
make dev

# 2. Fazer requisições
curl http://localhost:8080/health

# 3. Visualizar métricas no Jaeger UI
open http://localhost:16686

# 4. Ou configurar Grafana para consultar métricas
# Adicionar Jaeger como data source no Grafana
```

## Troubleshooting

### Métricas não aparecem
1. Verificar `SERVER_APP_OTEL_ENABLED=true`
2. Verificar conectividade com endpoint OTLP
3. Checar logs para erros de export

### Alto uso de memória
Aumentar intervalo de export:
```bash
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=30
```

### Métricas atrasadas
Diminuir intervalo de export:
```bash
SERVER_APP_OTEL_METRIC_EXPORT_INTERVAL=5
```

## Referências

- [Guia Completo de Métricas](./metrics-guide.md)
- [Guia de Observabilidade](./observability-guide.md)
- [OpenTelemetry Docs](https://opentelemetry.io/docs/instrumentation/go/manual/)
