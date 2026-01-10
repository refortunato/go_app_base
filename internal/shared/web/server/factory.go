package server

import "github.com/gin-gonic/gin"

// RouteSetupFunc defines a function that configures routes on a Gin router
// This allows generic server creation while keeping route definitions in infra layer
type RouteSetupFunc func(*gin.Engine)

// NewGinServerWithRoutes creates a new HTTP server with custom route setup
// The setupRoutes function is called to register application-specific routes
func NewGinServerWithRoutes(port string, setupRoutes RouteSetupFunc) *GinServer {
	if port == "" {
		port = "8080"
	}

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// Call the provided setup function to register routes
	if setupRoutes != nil {
		setupRoutes(router)
	}

	return NewGinServer(router, port)
}
