package server

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/shared/observability"
)

// RouteSetupFunc defines a function that configures routes on a Gin router
// This allows generic server creation while keeping route definitions in infra layer
type RouteSetupFunc func(*gin.Engine)

// NewGinServerWithRoutes creates a new HTTP server with custom route setup
// The setupRoutes function is called to register application-specific routes
func NewGinServerWithRoutes(port string, setupRoutes RouteSetupFunc, serviceName string, otelEnabled bool) *GinServer {
	if port == "" {
		port = "8080"
	}

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// Add OpenTelemetry tracing middleware if enabled
	if otelEnabled {
		router.Use(observability.TracingMiddleware(serviceName))
	}

	// Call the provided setup function to register routes
	if setupRoutes != nil {
		setupRoutes(router)
	}

	return NewGinServer(router, port)
}
