package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GinServer wraps http.Server for graceful shutdown
type GinServer struct {
	httpServer *http.Server
}

// Shutdown gracefully shuts down the server
func (s *GinServer) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

// Start starts the server and blocks until it's stopped
func (s *GinServer) Start() error {
	fmt.Printf("Starting HTTP server on %s\n", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// NewServer creates a new GinServer that can be gracefully shut down
func NewServer(webServerPort string) *GinServer {
	if webServerPort == "" {
		webServerPort = "8080"
	}

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()
	Handler(router)

	httpServer := &http.Server{
		Addr:           ":" + webServerPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &GinServer{
		httpServer: httpServer,
	}
}

// Start (deprecated) - kept for backward compatibility
// Use NewServer() instead for graceful shutdown support
func Start(webServerPort string) {
	server := NewServer(webServerPort)
	if err := server.Start(); err != nil {
		panic(err)
	}
}
