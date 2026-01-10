package dependencies

import (
	"database/sql"

	"github.com/refortunato/go_app_base/internal/core/application/repositories"
	"github.com/refortunato/go_app_base/internal/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/infra/config"
	infraLogger "github.com/refortunato/go_app_base/internal/infra/logger"
	infraRepositories "github.com/refortunato/go_app_base/internal/infra/repositories"
	"github.com/refortunato/go_app_base/internal/infra/web/controllers"
	"github.com/refortunato/go_app_base/internal/shared/logger"
)

// Repositories
var ExampleRepository repositories.ExampleRepository
var HealthRepository repositories.HealthRepository

// Logger
var Log logger.Logger

// UseCases
var GetExampleUseCase *usecases.GetExampleUseCase
var HealthCheckUseCase *usecases.HealthCheckUseCase

// Controllers
var ExampleController *controllers.ExampleController
var HealthController *controllers.HealthController

// Registrar a criação da dependência
func InitDependencies(db *sql.DB, config *config.Conf) error {
	// Logger
	Log = infraLogger.NewSlogLogger(config.ImageName, config.ImageVersion)
	logger.SetGlobalLogger(Log)
	logger.Info("Logger initialized successfully")

	// Repositories
	ExampleRepository = infraRepositories.NewExampleMySQLRepository(db)
	HealthRepository = infraRepositories.NewHealthMySQLRepository(db)

	// UseCases
	GetExampleUseCase = usecases.NewGetExampleUseCase(ExampleRepository)
	HealthCheckUseCase = usecases.NewHealthCheckUseCase(HealthRepository)

	// Controllers
	ExampleController = controllers.NewExampleController(*GetExampleUseCase)
	HealthController = controllers.NewHealthController(*HealthCheckUseCase)

	return nil
}
