package infra

import (
	"database/sql"

	"github.com/refortunato/go_app_base/internal/example/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/example/infra/repositories"
	"github.com/refortunato/go_app_base/internal/example/infra/web/controllers"
)

// ExampleModule encapsulates all dependencies for the example module
type ExampleModule struct {
	ExampleController *controllers.ExampleController
	GetExampleUseCase *usecases.GetExampleUseCase
}

// NewExampleModule creates and wires all dependencies for the example module
func NewExampleModule(db *sql.DB) *ExampleModule {
	// Repositories
	exampleRepository := repositories.NewExampleMySQLRepository(db)

	// Use Cases
	getExampleUseCase := usecases.NewGetExampleUseCase(exampleRepository)

	// Controllers
	exampleController := controllers.NewExampleController(*getExampleUseCase)

	return &ExampleModule{
		ExampleController: exampleController,
		GetExampleUseCase: getExampleUseCase,
	}
}
