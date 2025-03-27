package ezapp

// WireBundle is a struct that contains the components of the application.
//
// A WireBundle is returned by a WireFunc and contains:
// 1. A slice of Runnables that will be run by the EzApp
// 2. A CleanupFunc that will be called when the EzApp is done running
//
// The WireBundle is used by the Build function to create an EzApp.
// It serves as a container for the components that make up the application.
//
// Example:
//
//	func wireApp(cfg Config) (ezapp.WireBundle, error) {
//		server := NewServer(cfg.Port)
//		worker := NewWorker(cfg.WorkerConfig)
//
//		return ezapp.WireBundle{
//			Runnables: []ezapp.Runnable{server, worker},
//			CleanupFunc: func() error {
//				// Perform any cleanup operations
//				return nil
//			},
//		}, nil
//	}
type WireBundle struct {
	// Runnables is a slice of Runnable interfaces that will be run by the EzApp.
	// Each Runnable will be run in a separate goroutine.
	Runnables []Runnable

	// CleanupFunc is a function that will be called when the EzApp is done running.
	// It should release any resources that were allocated by the WireFunc.
	// If CleanupFunc returns an error, the EzApp will panic.
	CleanupFunc CleanupFunc
}
