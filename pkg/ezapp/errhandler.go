package ezapp

// ErrHandler is a function type that handles errors.
//
// An ErrHandler is called when a Runnable returns an error that it couldn't handle.
// It should attempt to handle the error and return nil if it was handled successfully,
// or return an error if it couldn't be handled.
//
// If an ErrHandler returns an error, the EzApp will cancel the context to initiate
// shutdown of all Runnables.
//
// Example:
//
//	func handleError(err error) error {
//		// Log the error
//		log.Printf("Error: %v", err)
//
//		// Check if it's a known error that we can handle
//		if errors.Is(err, SomeKnownError) {
//			// Handle the error and return nil to indicate it was handled
//			return nil
//		}
//
//		// Return the error to indicate it wasn't handled
//		return err
//	}
//
// An ErrHandler can be provided to the EzApp using the WithErrorHandler option
// when calling Build.
//
// The error handling in EzApp follows this chain:
// 1. If a Runnable returns an error, its HandleError method is called
// 2. If HandleError returns an error, the EzApp's error handler is called
// 3. If the EzApp's error handler returns an error, the context is canceled
type ErrHandler func(error) error
