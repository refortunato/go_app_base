package simple_module

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the simple_module (4-tier architecture)
func RegisterRoutes(router *gin.Engine, module *SimpleModule) {
	// Product routes
	router.GET("/products", func(ctx *gin.Context) {
		module.ProductController.ListProducts(context.NewGinContextAdapter(ctx))
	})

	router.GET("/products/:id", func(ctx *gin.Context) {
		module.ProductController.GetProduct(context.NewGinContextAdapter(ctx))
	})

	router.POST("/products", func(ctx *gin.Context) {
		module.ProductController.CreateProduct(context.NewGinContextAdapter(ctx))
	})

	router.PUT("/products/:id", func(ctx *gin.Context) {
		module.ProductController.UpdateProduct(context.NewGinContextAdapter(ctx))
	})

	router.DELETE("/products/:id", func(ctx *gin.Context) {
		module.ProductController.DeleteProduct(context.NewGinContextAdapter(ctx))
	})
}
