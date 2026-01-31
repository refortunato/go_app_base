package simple_module

import (
	"database/sql"

	"github.com/refortunato/go_app_base/internal/simple_module/controllers"
	"github.com/refortunato/go_app_base/internal/simple_module/repositories"
	"github.com/refortunato/go_app_base/internal/simple_module/services"
)

// SimpleModule holds all initialized dependencies for the simple_module (4-tier architecture)
// This module demonstrates a simpler architecture pattern for CRUD operations
type SimpleModule struct {
	ProductController *controllers.ProductController
	ProductService    *services.ProductService
}

// NewSimpleModule creates and wires all dependencies for the simple_module
func NewSimpleModule(db *sql.DB) *SimpleModule {
	// Step 1: Initialize repository
	productRepo := repositories.NewProductRepository(db)

	// Step 2: Initialize service (inject repository)
	productService := services.NewProductService(productRepo)

	// Step 3: Initialize controller (inject service)
	productController := controllers.NewProductController(productService)

	// Step 4: Return module with all dependencies wired
	return &SimpleModule{
		ProductController: productController,
		ProductService:    productService,
	}
}
