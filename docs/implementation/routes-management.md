# Routes Management Guide

This guide explains how to create and manage routes in the go_app_base project following the **Module Factory Pattern** and **Bounded Context** principles.

## Architecture Overview

Routes are managed at **two levels**:

1. **Module Level**: Each module defines its own routes in `internal/{module}/infra/web/routes.go`
2. **Application Level**: Central orchestrator in `internal/infra/web/register_routes.go` delegates to modules

This separation ensures **module independence** and prepares for potential microservice extraction.

---

## Creating Routes for a New Module

### Step 1: Create the Module's Routes File

**Location:** `internal/{module}/infra/web/routes.go`

**Template:**
```go
package web

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/{module}/infra/web/controllers"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the {module} module
func RegisterRoutes(router *gin.Engine, controller *controllers.{Module}Controller) {
	// Public routes
	router.GET("/{resource}/:id", func(ctx *gin.Context) {
		controller.Get{Resource}(context.NewGinContextAdapter(ctx))
	})

	router.POST("/{resource}", func(ctx *gin.Context) {
		controller.Create{Resource}(context.NewGinContextAdapter(ctx))
	})

	// Grouped routes (optional)
	resourceGroup := router.Group("/{resource}")
	{
		resourceGroup.PUT("/:id", func(ctx *gin.Context) {
			controller.Update{Resource}(context.NewGinContextAdapter(ctx))
		})

		resourceGroup.DELETE("/:id", func(ctx *gin.Context) {
			controller.Delete{Resource}(context.NewGinContextAdapter(ctx))
		})
	}
}
```

### Step 2: Update the Central Route Orchestrator

**Location:** `internal/infra/web/register_routes.go`

Add your module's route registration:

```go
package web

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/cmd/server/container"
	exampleWeb "github.com/refortunato/go_app_base/internal/example/infra/web"
	healthWeb "github.com/refortunato/go_app_base/internal/health/infra/web"
	// Add your module import
	yourModuleWeb "github.com/refortunato/go_app_base/internal/{module}/infra/web"
)

// RegisterRoutes is the main route orchestrator
func RegisterRoutes(c *container.Container) func(*gin.Engine) {
	return func(router *gin.Engine) {
		// Register routes for each module
		healthWeb.RegisterRoutes(router, c.HealthModule.HealthController)
		exampleWeb.RegisterRoutes(router, c.ExampleModule.ExampleController)
		
		// Add your module routes
		yourModuleWeb.RegisterRoutes(router, c.YourModule.YourController)
	}
}
```

---

## Route Patterns and Best Practices

### 1. REST Resource Routes

```go
// List resources
router.GET("/users", func(ctx *gin.Context) {
    controller.ListUsers(context.NewGinContextAdapter(ctx))
})

// Get single resource
router.GET("/users/:id", func(ctx *gin.Context) {
    controller.GetUser(context.NewGinContextAdapter(ctx))
})

// Create resource
router.POST("/users", func(ctx *gin.Context) {
    controller.CreateUser(context.NewGinContextAdapter(ctx))
})

// Update resource
router.PUT("/users/:id", func(ctx *gin.Context) {
    controller.UpdateUser(context.NewGinContextAdapter(ctx))
})

// Partial update
router.PATCH("/users/:id", func(ctx *gin.Context) {
    controller.PatchUser(context.NewGinContextAdapter(ctx))
})

// Delete resource
router.DELETE("/users/:id", func(ctx *gin.Context) {
    controller.DeleteUser(context.NewGinContextAdapter(ctx))
})
```

### 2. Nested Resources

```go
// Posts belonging to a user
router.GET("/users/:userId/posts", func(ctx *gin.Context) {
    controller.GetUserPosts(context.NewGinContextAdapter(ctx))
})

router.POST("/users/:userId/posts", func(ctx *gin.Context) {
    controller.CreateUserPost(context.NewGinContextAdapter(ctx))
})
```

### 3. Action Routes (Non-CRUD)

```go
// Custom actions
router.POST("/users/:id/activate", func(ctx *gin.Context) {
    controller.ActivateUser(context.NewGinContextAdapter(ctx))
})

router.POST("/orders/:id/cancel", func(ctx *gin.Context) {
    controller.CancelOrder(context.NewGinContextAdapter(ctx))
})
```

### 4. Route Groups with Middleware

```go
// Public routes
public := router.Group("/api/v1")
{
    public.GET("/health", func(ctx *gin.Context) {
        controller.HealthCheck(context.NewGinContextAdapter(ctx))
    })
}

// Protected routes (future: add auth middleware)
protected := router.Group("/api/v1")
// protected.Use(authMiddleware) // Add when implementing auth
{
    protected.GET("/users/:id", func(ctx *gin.Context) {
        controller.GetUser(context.NewGinContextAdapter(ctx))
    })
}
```

---

## Important: Use WebContext Abstraction

Always use `WebContext` interface instead of Gin directly:

**✅ CORRECT:**
```go
router.GET("/users/:id", func(ctx *gin.Context) {
    controller.GetUser(context.NewGinContextAdapter(ctx))
})
```

**❌ WRONG:**
```go
router.GET("/users/:id", controller.GetUser) // Don't pass gin.Context directly
```

**Why?**
- Keeps controllers framework-agnostic
- Easier to test (mock WebContext)
- Easier to switch web frameworks in the future

---

## Testing Module Routes

Create route tests in `internal/{module}/infra/web/routes_test.go`:

```go
package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/{module}/infra/web"
	"github.com/refortunato/go_app_base/internal/{module}/infra/web/controllers"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create mock controller
	controller := &controllers.MockController{}
	
	// Register routes
	web.RegisterRoutes(router, controller)
	
	// Test route exists
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/resource/123", nil)
	router.ServeHTTP(w, req)
	
	if w.Code == http.StatusNotFound {
		t.Error("Route not registered")
	}
}
```

---

## Route Organization Checklist

When adding routes to a module, ensure:

- ✅ Routes defined in `internal/{module}/infra/web/routes.go`
- ✅ Function signature: `RegisterRoutes(router *gin.Engine, controller *controllers.XController)`
- ✅ Use `context.NewGinContextAdapter(ctx)` for all handlers
- ✅ Module registered in `internal/infra/web/register_routes.go`
- ✅ RESTful naming conventions followed
- ✅ Routes logically grouped (if multiple resources)
- ✅ Controller methods exist for all routes

---

## Common Mistakes to Avoid

### ❌ Don't define routes in the container
```go
// WRONG: Don't do this in container.go
func New(db *sql.DB) *Container {
    router.GET("/users", ...) // NO!
}
```

### ❌ Don't couple modules through routes
```go
// WRONG: Don't import other modules
import "github.com/refortunato/go_app_base/internal/other_module/..."
```

### ❌ Don't bypass WebContext
```go
// WRONG: Don't use gin.Context directly in controller
func (c *Controller) GetUser(ctx *gin.Context) { // NO!
```

---

## Migration from Old Structure

If migrating from centralized routes:

1. Create `routes.go` in each module's `infra/web/` directory
2. Move module-specific routes from central file to module's `routes.go`
3. Update central orchestrator to call module's `RegisterRoutes()`
4. Test each module independently

---

## Summary

**Module Routes** (`internal/{module}/infra/web/routes.go`):
- Define routes specific to the module
- Export `RegisterRoutes(router, controller)` function
- Use WebContext abstraction
- Keep module independent

**Central Orchestrator** (`internal/infra/web/register_routes.go`):
- Calls each module's `RegisterRoutes()`
- Does NOT define routes itself
- Acts as a delegation layer
- Receives container to access module controllers

This pattern ensures **bounded context independence** and makes it easy to extract modules into separate microservices in the future.
