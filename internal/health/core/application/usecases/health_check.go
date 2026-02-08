package usecases

import (
	"context"

	"github.com/refortunato/go_app_base/internal/health/core/application/repositories"
	"github.com/refortunato/go_app_base/internal/shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type HealthCheckOutputDTO struct {
	Status string `json:"status"`
}

type HealthCheckUseCase struct {
	healthRepository repositories.HealthRepository
	metrics          *observability.CustomMetrics
	healthCounter    metric.Int64Counter
}

func NewHealthCheckUseCase(healthRepository repositories.HealthRepository) *HealthCheckUseCase {
	metrics := observability.NewCustomMetrics("health_module")

	// Create counter for health checks (reuse across all calls)
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

	// Record health check metric (non-blocking)
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

	output := &HealthCheckOutputDTO{
		Status: "OK",
	}

	return output, nil
}
