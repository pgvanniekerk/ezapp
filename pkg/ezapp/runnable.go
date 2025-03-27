package ezapp

import (
	"context"
)

// Runnable is an interface for components that can be run by the EzApp.
//
// A Runnable represents a long-running component of the application, such as
// an HTTP server, a background worker, or a message consumer. The EzApp runs
// each Runnable in a separate goroutine and coordinates their lifecycle.
//
// Implementers of this interface should:
// 1. Implement Run to start the component and keep it running until the context is canceled
// 2. Implement HandleError to handle errors that occur during the run
//
// Example implementation:
//
//	type Server struct {
//		server *http.Server
//	}
//
//	func (s *Server) Run(ctx context.Context) error {
//		// Start the server in a goroutine
//		errCh := make(chan error, 1)
//		go func() {
//			errCh <- s.server.ListenAndServe()
//		}()
//
//		// Wait for context cancellation or server error
//		select {
//		case <-ctx.Done():
//			// Context was canceled, shut down gracefully
//			return s.server.Shutdown(context.Background())
//		case err := <-errCh:
//			// Server encountered an error
//			return err
//		}
//	}
//
//	func (s *Server) HandleError(err error) error {
//		// Log the error and return nil to indicate it was handled
//		log.Printf("Server error: %v", err)
//		return nil
//	}
type Runnable interface {
	// Run starts the component and keeps it running until the context is canceled.
	// It should return an error if the component fails to start or encounters an error during execution.
	// If the context is canceled, Run should perform a graceful shutdown and return context.Canceled.
	Run(context.Context) error

	// HandleError handles errors that occur during the run.
	// It should return nil if the error was handled successfully, or an error if it couldn't be handled.
	// If HandleError returns an error, the EzApp will attempt to handle it with its error handler.
	HandleError(error) error
}
