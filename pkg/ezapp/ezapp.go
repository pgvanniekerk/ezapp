package ezapp

// EzApp is the core interface that represents a runnable application.
// Any type that implements this interface can be used with the ezapp framework.
//
// The ezapp framework provides utilities for building, configuring, and running
// applications that implement this interface. Typically, applications will be
// created using the wire.App function and run using the ezapp.Run function.
//
// Example:
//
//	func BuildFunc(ctx context.Context, cfg Config) (ezapp.EzApp, error) {
//	    return wire.App(
//	        wire.Runnables(myRunnable),
//	        wire.WithAppShutdownTimeout(15*time.Second),
//	    )
//	}
//
//	func main() {
//	    ezapp.Run(BuildFunc)
//	}
type EzApp interface {
	// Run starts the application and blocks until the application exits.
	// This method is called by ezapp.Run and should handle the entire
	// application lifecycle, including startup, running, and shutdown.
	Run()
}
