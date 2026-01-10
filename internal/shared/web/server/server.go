package server

import "context"

// Server represents a service that can be started and gracefully shut down
// This interface can be implemented by HTTP servers, gRPC servers, message consumers, etc.
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}
