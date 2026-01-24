# Dependency Management Guide

This guide explains how to manage dependencies in the go_app_base project using the **Module Factory Pattern** and **Composition Root** principles.

## Architecture Overview

Dependency injection happens at **three levels**:

1. **Module Level**: Each module wires its own dependencies in `internal/{module}/infra/module.go`
2. **Container Level**: Application container orchestrates modules in `cmd/server/container/container.go`
3. **Main Level**: Entry point creates container and passes it to server in `cmd/server/main.go`

This separation ensures **low coupling**, **high cohesion**, and **testability**.

---

## Dependency Flow

```
main.go
  ↓ creates
container.Container
  ↓ initializes
{Module}Module (via New{Module}Module factory)
  ↓ wires
Repository → UseCase → Controller
```

---

## Creating Dependencies for a New Module

### Step 1: Create the Module Factory

**Location:** `internal/{module}/infra/module.go`

**Template:**
```go
package infra

import (
	"database/sql"

	"github.com/refortunato/go_app_base/internal/{module}/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/{module}/infra/repositories"
	"github.com/refortunato/go_app_base/internal/{module}/infra/web/controllers"
)

// {Module}Module holds all initialized dependencies for the {module} module
type {Module}Module struct {
	// Controllers (for web layer)
	{Resource}Controller *controllers.{Resource}Controller

	// Use Cases (for direct access if needed)
	Get{Resource}UseCase    *usecases.Get{Resource}UseCase
	Create{Resource}UseCase *usecases.Create{Resource}UseCase
	Update{Resource}UseCase *usecases.Update{Resource}UseCase
	Delete{Resource}UseCase *usecases.Delete{Resource}UseCase
}

// New{Module}Module creates and wires all dependencies for the {module} module
func New{Module}Module(db *sql.DB) *{Module}Module {
	// Step 1: Initialize repositories (infrastructure → application boundary)
	resourceRepo := repositories.New{Resource}MySQLRepository(db)

	// Step 2: Initialize use cases (inject repository dependencies)
	getResourceUseCase := usecases.NewGet{Resource}UseCase(resourceRepo)
	createResourceUseCase := usecases.NewCreate{Resource}UseCase(resourceRepo)
	updateResourceUseCase := usecases.NewUpdate{Resource}UseCase(resourceRepo)
	deleteResourceUseCase := usecases.NewDelete{Resource}UseCase(resourceRepo)

	// Step 3: Initialize controllers (inject use case dependencies)
	resourceController := controllers.New{Resource}Controller(
		*getResourceUseCase,
		*createResourceUseCase,
		*updateResourceUseCase,
		*deleteResourceUseCase,
	)

	// Step 4: Return module with all dependencies wired
	return &{Module}Module{
		{Resource}Controller:    resourceController,
		Get{Resource}UseCase:    getResourceUseCase,
		Create{Resource}UseCase: createResourceUseCase,
		Update{Resource}UseCase: updateResourceUseCase,
		Delete{Resource}UseCase: deleteResourceUseCase,
	}
}
```

### Step 2: Register Module in Container

**Location:** `cmd/server/container/container.go`

Add your module to the container:

```go
package container

import (
	"database/sql"

	"github.com/refortunato/go_app_base/configs"
	exampleInfra "github.com/refortunato/go_app_base/internal/example/infra"
	healthInfra "github.com/refortunato/go_app_base/internal/health/infra"
	"github.com/refortunato/go_app_base/internal/shared/logger"
	
	// Add your module import
	yourModuleInfra "github.com/refortunato/go_app_base/internal/{module}/infra"
)

// Container holds all initialized modules (Composition Root)
type Container struct {
	Logger        logger.Logger
	ExampleModule *exampleInfra.ExampleModule
	HealthModule  *healthInfra.HealthModule
	
	// Add your module
	YourModule    *yourModuleInfra.{Module}Module
}

// New creates and initializes the application container with all modules
func New(db *sql.DB, cfg *configs.Conf) (*Container, error) {
	// Initialize shared logger
	appLogger, err := logger.NewSlogLogger()
	if err != nil {
		return nil, err
	}

	// Initialize modules (each module wires its own dependencies)
	exampleModule := exampleInfra.NewExampleModule(db)
	healthModule := healthInfra.NewHealthModule(db)
	
	// Initialize your module
	yourModule := yourModuleInfra.New{Module}Module(db)

	return &Container{
		Logger:        appLogger,
		ExampleModule: exampleModule,
		HealthModule:  healthModule,
		YourModule:    yourModule,
	}, nil
}
```

---

## Dependency Types and Injection Points

### 1. Repository Dependencies

**Interface location:** `internal/{module}/core/application/repositories/{resource}_repository.go`

**Implementation location:** `internal/{module}/infra/repositories/{resource}_mysql_repository.go`

**Injection point:** Use case constructor

```go
// Interface (core/application/repositories)
type ExampleRepository interface {
	Save(example *entities.Example) error
	FindById(id string) (*entities.Example, error)
	Update(example *entities.Example) error
	Delete(id string) error
}

// Implementation (infra/repositories)
type ExampleMySQLRepository struct {
	db *sql.DB
}

func NewExampleMySQLRepository(db *sql.DB) *ExampleMySQLRepository {
	return &ExampleMySQLRepository{db: db}
}

// Injection in module.go
exampleRepo := repositories.NewExampleMySQLRepository(db)
useCase := usecases.NewGetExampleUseCase(exampleRepo)
```

### 2. Use Case Dependencies

**Location:** `internal/{module}/core/application/usecases/{action}_{resource}.go`

**Injection point:** Controller constructor

```go
// Use Case (core/application/usecases)
type GetExampleUseCase struct {
	repository repositories.ExampleRepository
}

func NewGetExampleUseCase(repo repositories.ExampleRepository) *GetExampleUseCase {
	return &GetExampleUseCase{repository: repo}
}

// Injection in module.go
getExampleUseCase := usecases.NewGetExampleUseCase(exampleRepo)
controller := controllers.NewExampleController(*getExampleUseCase)
```

### 3. Controller Dependencies

**Location:** `internal/{module}/infra/web/controllers/{resource}_controller.go`

**Injection point:** Module factory → used by routes

```go
// Controller (infra/web/controllers)
type ExampleController struct {
	getExampleUseCase GetExampleUseCase
}

func NewExampleController(getUseCase GetExampleUseCase) *ExampleController {
	return &ExampleController{getExampleUseCase: getUseCase}
}

// Injection in module.go
controller := controllers.NewExampleController(*getExampleUseCase)

// Usage in routes.go
web.RegisterRoutes(router, c.ExampleModule.ExampleController)
```

### 4. Multiple Repository Dependencies

When a use case needs multiple repositories:

```go
// Use Case with multiple dependencies
type CreateOrderUseCase struct {
	orderRepository   repositories.OrderRepository
	productRepository repositories.ProductRepository
	userRepository    repositories.UserRepository
}

func NewCreateOrderUseCase(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	userRepo repositories.UserRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepository:   orderRepo,
		productRepository: productRepo,
		userRepository:    userRepo,
	}
}

// Injection in module.go
orderRepo := repositories.NewOrderMySQLRepository(db)
productRepo := repositories.NewProductMySQLRepository(db)
userRepo := repositories.NewUserMySQLRepository(db)

createOrderUseCase := usecases.NewCreateOrderUseCase(
	orderRepo,
	productRepo,
	userRepo,
)
```

### 5. Domain Service Dependencies

When use cases need domain services:

```go
// Domain Service (core/domain/services)
type PricingService struct{}

func NewPricingService() *PricingService {
	return &PricingService{}
}

// Use Case with domain service
type CreateOrderUseCase struct {
	orderRepository repositories.OrderRepository
	pricingService  *services.PricingService
}

func NewCreateOrderUseCase(
	orderRepo repositories.OrderRepository,
	pricingService *services.PricingService,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepository: orderRepo,
		pricingService:  pricingService,
	}
}

// Injection in module.go
pricingService := services.NewPricingService()
orderRepo := repositories.NewOrderMySQLRepository(db)
createOrderUseCase := usecases.NewCreateOrderUseCase(orderRepo, pricingService)
```

---

## Adding Dependencies to Existing Modules

### Step 1: Create the New Dependency

Example: Adding a new use case to an existing module

**File:** `internal/example/core/application/usecases/update_example.go`

```go
package usecases

import (
	"github.com/refortunato/go_app_base/internal/example/core/application/repositories"
	"github.com/refortunato/go_app_base/internal/example/core/domain/entities"
)

type UpdateExampleUseCase struct {
	repository repositories.ExampleRepository
}

func NewUpdateExampleUseCase(repo repositories.ExampleRepository) *UpdateExampleUseCase {
	return &UpdateExampleUseCase{repository: repo}
}

func (uc *UpdateExampleUseCase) Execute(input UpdateExampleInputDTO) error {
	// Implementation
}
```

### Step 2: Update Module Factory

**File:** `internal/example/infra/module.go`

```go
type ExampleModule struct {
	ExampleController *controllers.ExampleController
	GetExampleUseCase *usecases.GetExampleUseCase
	UpdateExampleUseCase *usecases.UpdateExampleUseCase // Add new dependency
}

func NewExampleModule(db *sql.DB) *ExampleModule {
	exampleRepo := repositories.NewExampleMySQLRepository(db)
	getExampleUseCase := usecases.NewGetExampleUseCase(exampleRepo)
	updateExampleUseCase := usecases.NewUpdateExampleUseCase(exampleRepo) // Wire new dependency

	exampleController := controllers.NewExampleController(
		*getExampleUseCase,
		*updateExampleUseCase, // Pass to controller if needed
	)

	return &ExampleModule{
		ExampleController:    exampleController,
		GetExampleUseCase:    getExampleUseCase,
		UpdateExampleUseCase: updateExampleUseCase, // Add to module
	}
}
```

### Step 3: Update Controller (if needed)

**File:** `internal/example/infra/web/controllers/example_controller.go`

```go
type ExampleController struct {
	getExampleUseCase    usecases.GetExampleUseCase
	updateExampleUseCase usecases.UpdateExampleUseCase // Add new dependency
}

func NewExampleController(
	getUseCase usecases.GetExampleUseCase,
	updateUseCase usecases.UpdateExampleUseCase, // Receive new dependency
) *ExampleController {
	return &ExampleController{
		getExampleUseCase:    getUseCase,
		updateExampleUseCase: updateUseCase,
	}
}
```

---

## Shared Dependencies

For dependencies used across multiple modules (e.g., logger, config):

### Option 1: Pass via Module Factory

```go
func NewExampleModule(db *sql.DB, logger logger.Logger) *ExampleModule {
	exampleRepo := repositories.NewExampleMySQLRepository(db, logger)
	// ...
}

// In container.go
exampleModule := exampleInfra.NewExampleModule(db, appLogger)
```

### Option 2: Pass to Specific Components

```go
// In module.go
func NewExampleModule(db *sql.DB, logger logger.Logger) *ExampleModule {
	exampleRepo := repositories.NewExampleMySQLRepository(db)
	
	// Pass logger only where needed
	getExampleUseCase := usecases.NewGetExampleUseCase(exampleRepo, logger)
	
	exampleController := controllers.NewExampleController(*getExampleUseCase)
	
	return &ExampleModule{
		ExampleController: exampleController,
		GetExampleUseCase: getExampleUseCase,
	}
}
```

---

## Testing with Dependency Injection

### Unit Test Example (Use Case)

```go
package usecases_test

import (
	"testing"

	"github.com/refortunato/go_app_base/internal/example/core/application/usecases"
	"github.com/refortunato/go_app_base/internal/example/core/application/repositories"
)

// Mock repository
type MockExampleRepository struct{}

func (m *MockExampleRepository) FindById(id string) (*entities.Example, error) {
	return &entities.Example{ID: id}, nil
}

func TestGetExampleUseCase(t *testing.T) {
	// Arrange: inject mock repository
	mockRepo := &MockExampleRepository{}
	useCase := usecases.NewGetExampleUseCase(mockRepo)
	
	// Act
	input := usecases.GetExampleInputDTO{ID: "123"}
	result, err := useCase.Execute(input)
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.ID != "123" {
		t.Errorf("Expected ID 123, got %s", result.ID)
	}
}
```

### Integration Test Example (Module Factory)

```go
package infra_test

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/refortunato/go_app_base/internal/example/infra"
)

func TestNewExampleModule(t *testing.T) {
	// Setup test database
	db, err := sql.Open("mysql", "test_dsn")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	
	// Create module
	module := infra.NewExampleModule(db)
	
	// Verify all dependencies are wired
	if module.ExampleController == nil {
		t.Error("ExampleController not initialized")
	}
	if module.GetExampleUseCase == nil {
		t.Error("GetExampleUseCase not initialized")
	}
}
```

---

## Dependency Management Checklist

When adding new dependencies:

- ✅ Define repository interfaces in `core/application/repositories`
- ✅ Implement repositories in `infra/repositories`
- ✅ Create use cases in `core/application/usecases`
- ✅ Use constructor injection (no global variables)
- ✅ Wire dependencies in module factory (`infra/module.go`)
- ✅ Update module struct to expose needed dependencies
- ✅ Register module in container if it's new
- ✅ Dependencies flow: Repository → UseCase → Controller
- ✅ Each dependency has a constructor function (`NewXxx`)

---

## Common Mistakes to Avoid

### ❌ Don't use global variables for dependencies

```go
// WRONG
var globalDB *sql.DB

func GetExampleUseCase() {
	globalDB.Query(...) // Don't use globals
}
```

### ❌ Don't wire dependencies in main.go

```go
// WRONG: Don't do detailed wiring in main.go
func main() {
	repo := repositories.NewExampleRepository(db)
	useCase := usecases.NewGetExampleUseCase(repo)
	controller := controllers.NewExampleController(useCase)
	// This should be in module.go!
}
```

### ❌ Don't create dependencies inside use cases

```go
// WRONG: Don't instantiate dependencies internally
func (uc *GetExampleUseCase) Execute(input InputDTO) error {
	repo := repositories.NewExampleRepository(db) // NO! Inject via constructor
}
```

### ❌ Don't expose module internals through container

```go
// WRONG: Don't expose individual use cases/repos if not needed
type Container struct {
	ExampleRepository repositories.ExampleRepository // NO! Keep internal
	GetExampleUseCase *usecases.GetExampleUseCase    // Only if needed externally
}

// CORRECT: Expose only what's needed (usually just controllers)
type Container struct {
	ExampleModule *exampleInfra.ExampleModule // Module encapsulates internals
}
```

---

## Benefits of This Pattern

1. **Testability**: Easy to inject mocks in tests
2. **Flexibility**: Easy to swap implementations
3. **Independence**: Modules don't depend on each other
4. **Clarity**: Dependency graph is explicit and visible
5. **Maintainability**: Each module manages its own dependencies
6. **Microservice-Ready**: Modules can be extracted independently

---

## Summary

**Module Factory** (`internal/{module}/infra/module.go`):
- Wires all dependencies for one module
- Returns a module struct with initialized components
- Keeps wiring logic close to implementation
- Each module is self-contained

**Container** (`cmd/server/container/container.go`):
- Orchestrates modules only
- Does NOT wire individual dependencies
- Passes shared resources (db, config) to modules
- Thin composition root

**Main** (`cmd/server/main.go`):
- Creates DB connection
- Loads configuration
- Creates container
- Starts server
- **Does NOT wire dependencies**

This pattern ensures **proper separation of concerns** and makes the codebase **easy to maintain and test**.
