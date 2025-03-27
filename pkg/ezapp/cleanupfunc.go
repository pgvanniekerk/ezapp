package ezapp

// CleanupFunc is a function type that performs cleanup operations.
//
// A CleanupFunc is called when the EzApp is done running, either because all
// Runnables have completed or because the application is shutting down due to
// a signal or an error. It should release any resources that were allocated
// by the application, such as database connections, file handles, or network
// connections.
//
// If a CleanupFunc returns an error, the EzApp will panic.
//
// Example:
//
//	func cleanup() error {
//		// Close database connection
//		if err := db.Close(); err != nil {
//			return fmt.Errorf("failed to close database connection: %w", err)
//		}
//
//		// Close file handles
//		if err := file.Close(); err != nil {
//			return fmt.Errorf("failed to close file: %w", err)
//		}
//
//		return nil
//	}
//
// A CleanupFunc can be provided to the EzApp in two ways:
// 1. As part of the WireBundle returned by the WireFunc
// 2. Using the WithCleanupFunc option when calling Build
//
// If both are provided, the one from the WithCleanupFunc option takes precedence.
type CleanupFunc func() error
