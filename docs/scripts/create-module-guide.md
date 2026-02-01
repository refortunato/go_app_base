# Create Module Script Guide

The `create-module.sh` script automates the creation of new bounded context modules in your Go application, supporting both **DDD (Clean Architecture)** and **4-tier (simplified)** architectures.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Usage](#usage)
- [Architecture Options](#architecture-options)
- [What Gets Created](#what-gets-created)
- [Step-by-Step Example](#step-by-step-example)
- [Next Steps](#next-steps)

## Overview

This script helps you:
- Create a complete module structure following architectural best practices
- Automatically wire dependencies in the container
- Register routes in the application
- Choose between two architecture styles based on your needs

## Prerequisites

- Go 1.25+ installed
- Project initialized with `go.mod`
- Write permissions in the project directory

## Usage

From the project root, run:

```bash
./scripts/create-module.sh
```

The script will guide you through an interactive process:

1. **Module Name**: Enter a descriptive name (e.g., `user`, `order`, `payment`)
2. **Architecture Type**: Choose between:
   - `1` - 4-tier (simplified) - For simple CRUD operations
   - `2` - DDD (Clean Architecture) - For complex business logic

## Architecture Options

### Option 1: 4-Tier (Simplified)

**Use when:**
- Simple CRUD operations
- Basic validations
- Straightforward data access
- Fast development needed

**Structure:**
```
internal/{module}/
├── models/           # Data structures
├── repositories/     # Database access
├── services/         # Business logic
├── controllers/      # HTTP handlers
├── routes.go         # Route definitions
└── module.go         # Dependency wiring
```

**Example modules:** `simple_module/`

### Option 2: DDD (Clean Architecture)

**Use when:**
- Complex business rules
- Multiple aggregates
- Rich domain logic
- Domain events needed
- Long-term maintainability critical

**Structure:**
```
internal/{module}/
├── core/
│   ├── application/
│   │   ├── repositories/    # Repository interfaces
│   │   ├── usecases/        # Use cases
│   │   └── services/        # Application services
│   └── domain/
│       ├── entities/        # Domain entities
│       ├── valueobjects/    # Value objects
│       └── services/        # Domain services
└── infra/
    ├── repositories/        # Repository implementations
    ├── web/
    │   ├── controllers/     # HTTP controllers
    │   └── routes.go        # Route definitions
    └── module.go            # Dependency wiring
```

**Example modules:** `example/`, `health/`, `payment/`

## What Gets Created

### For Both Architectures:

1. **Complete directory structure** based on chosen architecture
2. **module.go** - Factory for dependency injection with TODO comments
3. **routes.go** - Route registration file with examples
4. **Automatic updates to:**
   - `cmd/server/container/container.go` - Adds module to container
   - `internal/infra/web/register_routes.go` - Registers module routes

### Naming Conventions:

- **Module name**: lowercase with underscores (e.g., `user_management`)
- **Go types**: PascalCase (e.g., `UserManagementModule`)
- **Package names**: lowercase (e.g., `package user_management`)

## Step-by-Step Example

### Example 1: Creating a Simple User Module (4-tier)

```bash
$ ./scripts/create-module.sh

ℹ Project module path: github.com/company/myapp

ℹ Enter the module name (e.g., user, order, payment):
user

ℹ Module name: user

ℹ Select architecture type:
1) 4-tier (simplified) - For simple CRUD operations
2) DDD (Clean Architecture) - For complex business logic
Enter your choice (1 or 2): 1

ℹ Creating module 'user'...
ℹ Creating 4-tier (simplified) architecture...
✓ Created directory structure
✓ Created module.go
✓ Created routes.go
✓ Added import to container.go
✓ Added field to Container struct
✓ Added module initialization to New function
✓ Added field to Container return statement
✓ Added import to register_routes.go
✓ Added route registration to register_routes.go

✓ Module 'user' created successfully!

ℹ Next steps:
  1. Implement your domain entities, repositories, services, and controllers
  2. Define your routes in the routes.go file
  3. Wire your dependencies in the module.go file
  4. Test your module by running the application

ℹ Module location: internal/user/
ℹ Architecture: 4-tier (simplified)
  - models/       - Data structures
  - repositories/ - Database access
  - services/     - Business logic
  - controllers/  - HTTP handlers
```

### Example 2: Creating a Payment Module (DDD)

```bash
$ ./scripts/create-module.sh

ℹ Enter the module name (e.g., user, order, payment):
payment

ℹ Select architecture type:
1) 4-tier (simplified) - For simple CRUD operations
2) DDD (Clean Architecture) - For complex business logic
Enter your choice (1 or 2): 2

ℹ Creating module 'payment'...
ℹ Creating DDD (Clean Architecture) structure...
✓ Created directory structure
✓ Created infra/module.go
✓ Created infra/web/routes.go
...

ℹ Module location: internal/payment/
ℹ Architecture: DDD (Clean Architecture)
  - core/application/repositories/ - Repository interfaces
  - core/application/usecases/     - Use cases
  - core/domain/entities/          - Domain entities
  - infra/repositories/            - Repository implementations
  - infra/web/controllers/         - HTTP controllers
```

## Next Steps

After creating a module, you need to:

### 1. Implement Your Components

**For 4-tier:**
- Define models in `models/`
- Implement repositories in `repositories/`
- Create services in `services/`
- Add controllers in `controllers/`

**For DDD:**
- Define domain entities in `core/domain/entities/`
- Create repository interfaces in `core/application/repositories/`
- Implement use cases in `core/application/usecases/`
- Implement repositories in `infra/repositories/`
- Add controllers in `infra/web/controllers/`

### 2. Wire Dependencies in module.go

Open the generated `module.go` file and follow the TODO comments:

**4-tier example:**
```go
func NewUserModule(db *sql.DB) *UserModule {
    // Step 1: Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    
    // Step 2: Initialize services (inject repositories)
    userService := services.NewUserService(userRepo)
    
    // Step 3: Initialize controllers (inject services)
    userController := controllers.NewUserController(userService)
    
    // Step 4: Return module with all dependencies wired
    return &UserModule{
        UserController: userController,
    }
}
```

**DDD example:**
```go
func NewPaymentModule(db *sql.DB) *PaymentModule {
    // Step 1: Initialize repositories
    paymentRepo := repositories.NewPaymentMySQLRepository(db)
    
    // Step 2: Initialize use cases (inject repositories)
    createPaymentUC := usecases.NewCreatePaymentUseCase(paymentRepo)
    getPaymentUC := usecases.NewGetPaymentUseCase(paymentRepo)
    
    // Step 3: Initialize controllers (inject use cases)
    paymentController := controllers.NewPaymentController(
        *createPaymentUC,
        *getPaymentUC,
    )
    
    // Step 4: Return module with all dependencies wired
    return &PaymentModule{
        PaymentController: paymentController,
    }
}
```

### 3. Define Routes in routes.go

**4-tier example:**
```go
func RegisterRoutes(router *gin.Engine, module *UserModule) {
    router.POST("/users", func(ctx *gin.Context) {
        module.UserController.Create(context.NewGinContextAdapter(ctx))
    })
    
    router.GET("/users/:id", func(ctx *gin.Context) {
        module.UserController.Get(context.NewGinContextAdapter(ctx))
    })
}
```

**DDD example:**
```go
func RegisterRoutes(router *gin.Engine, module *infra.PaymentModule) {
    router.POST("/payments", func(ctx *gin.Context) {
        module.PaymentController.Create(context.NewGinContextAdapter(ctx))
    })
    
    router.GET("/payments/:id", func(ctx *gin.Context) {
        module.PaymentController.Get(context.NewGinContextAdapter(ctx))
    })
}
```

### 4. Use the Entity Script (Optional)

To quickly scaffold entities with full CRUD operations, use:

```bash
./scripts/create-entity.sh
```

See [Create Entity Guide](./create-entity-guide.md) for details.

### 5. Test Your Module

```bash
# Run the application
go run cmd/server/main.go

# Test your endpoints
curl http://localhost:8080/users
```

## Module Independence

Each module is designed to be:
- **Self-contained**: All dependencies are within the module
- **Loosely coupled**: Depends only on `shared/` and `configs/`
- **Extractable**: Can be moved to a separate microservice without changes

## Best Practices

1. **Choose the right architecture**: 4-tier for simple features, DDD for complex domains
2. **Keep modules independent**: Don't import other modules
3. **Wire all dependencies**: Use the module factory pattern
4. **Register routes properly**: Each module registers its own routes
5. **Follow naming conventions**: Consistent names across the project

## Troubleshooting

### Module Already Exists
```
✗ Module 'user' already exists at internal/user
```
**Solution**: Choose a different name or delete the existing module

### Import Not Added
```
⚠ Import already exists in container.go
```
**Solution**: This is just a warning. The script detected an existing import.

### Container Not Updated
If the container wasn't updated automatically, manually add:

1. Import: `userModule "github.com/company/myapp/internal/user"`
2. Field: `UserModule *userModule.UserModule`
3. Init: `userModule := userModule.NewUserModule(db)`
4. Return: `UserModule: userModule,`

## See Also

- [Create Entity Guide](./create-entity-guide.md) - Scaffold complete CRUD entities
- [Architecture Decision](../implementation/routes-management.md) - Understanding the routing approach
- [Dependency Management](../implementation/dependency-management.md) - How DI works in this project
