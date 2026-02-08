package usecases

import (
	"context"
	"time"

	"github.com/refortunato/go_app_base/internal/example/core/application/repositories"
	"github.com/refortunato/go_app_base/internal/example/core/domain/entities"
	"github.com/refortunato/go_app_base/internal/shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Example use case with comprehensive metrics instrumentation
// Demonstrates different metric types and best practices

type CreateExampleMetricsDemo struct {
	repository repositories.ExampleRepository

	// Metrics instruments (created once, reused many times)
	metrics          *observability.CustomMetrics
	creationCounter  metric.Int64Counter       // Total created
	creationDuration metric.Float64Histogram   // Time to create
	activeCreations  metric.Int64UpDownCounter // In-progress operations
}

func NewCreateExampleMetricsDemo(repo repositories.ExampleRepository) *CreateExampleMetricsDemo {
	metrics := observability.NewCustomMetrics("example_module")

	// Initialize all metric instruments upfront (efficient reuse)
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
		repository:       repo,
		metrics:          metrics,
		creationCounter:  creationCounter,
		creationDuration: creationDuration,
		activeCreations:  activeCreations,
	}
}

func (uc *CreateExampleMetricsDemo) Execute(ctx context.Context, name string) (*entities.Example, error) {
	// Track active operations (increment)
	uc.activeCreations.Add(ctx, 1)
	defer func() {
		// Decrement on completion (non-blocking)
		uc.activeCreations.Add(ctx, -1)
	}()

	// Measure operation duration
	start := time.Now()

	// Create entity
	example, err := entities.NewExample(name)
	if err != nil {
		// Record failure metric
		uc.creationCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("status", "validation_error"),
			),
		)
		return nil, err
	}

	// Save to repository
	err = uc.repository.Save(example)
	if err != nil {
		// Record failure metric
		uc.creationCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("status", "repository_error"),
			),
		)
		return nil, err
	}

	// Calculate duration
	duration := float64(time.Since(start).Milliseconds())

	// Record success metrics (all non-blocking)
	uc.creationCounter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("status", "success"),
		),
	)

	uc.creationDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.String("operation", "create"),
		),
	)

	return example, nil
}

// Example of using async gauge for monitoring repository state
func (uc *CreateExampleMetricsDemo) RegisterGaugeMetrics(repo repositories.ExampleRepository) error {
	// This callback is executed asynchronously and periodically
	// Should NOT block or perform heavy operations
	return uc.metrics.Gauge(
		"examples.repository.size",
		"Approximate number of examples in repository",
		"{example}",
		func(ctx context.Context, observer metric.Int64Observer) error {
			// In real implementation, this would query a cached count
			// NOT a full COUNT(*) query (too expensive for periodic callback)
			count := int64(0) // Replace with cached value
			observer.Observe(count)
			return nil
		},
	)
}
