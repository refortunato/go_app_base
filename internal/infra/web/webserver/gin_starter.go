package webserver

import (
	"github.com/gin-gonic/gin"
)

func Start(webServerPort string) {
	if webServerPort == "" {
		webServerPort = "8080"
	}
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	Handler(router)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	router.Run(":" + webServerPort)
}
