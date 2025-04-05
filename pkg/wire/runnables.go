package wire

import "github.com/pgvanniekerk/ezapp/internal/app"

// Runnables creates a function that returns the provided runnable components.
// This function is used with the App function to provide the runnable components
// that will be managed by the application.
//
// Parameters:
//   - runnables: A variadic list of Runnable components that implement the app.Runnable interface.
//     These are typically structs that embed ezapp.Runnable and override the Run and Stop methods.
//
// Returns:
//   - A function that returns the provided runnables as a slice.
//     This function is used as the first argument to the App function.
//
// Example:
//
//	// Create runnable components
//	dbRunnable := myapp.NewDBRunnable(db)
//	httpRunnable := myapp.NewHTTPRunnable(router)
//
//	// Create the application with the runnables
//	app, err := wire.App(
//	    wire.Runnables(dbRunnable, httpRunnable),
//	    wire.WithAppShutdownTimeout(15*time.Second),
//	)
//
// Note: The App function expects at least one runnable component. If no runnables
// are provided, the application will still be created but won't do anything useful.
func Runnables(runnables ...app.Runnable) func() []app.Runnable {
	return func() []app.Runnable {
		return runnables
	}
}
