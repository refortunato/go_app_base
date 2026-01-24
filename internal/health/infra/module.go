package infra

import (
	"database/sql"

	"github.com/refortunato/go_app_base/internal/health/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/health/infra/repositories"
	"github.com/refortunato/go_app_base/internal/health/infra/web/controllers"
)

// HealthModule encapsulates all dependencies for the health module
type HealthModule struct {
	HealthController   *controllers.HealthController
	HealthCheckUseCase *usecases.HealthCheckUseCase
}

// NewHealthModule creates and wires all dependencies for the health module
func NewHealthModule(db *sql.DB) *HealthModule {
	// Repositories
	healthRepository := repositories.NewHealthMySQLRepository(db)

	// Use Cases
	healthCheckUseCase := usecases.NewHealthCheckUseCase(healthRepository)

	// Controllers
	healthController := controllers.NewHealthController(*healthCheckUseCase)

	return &HealthModule{
		HealthController:   healthController,
		HealthCheckUseCase: healthCheckUseCase,
	}
}
