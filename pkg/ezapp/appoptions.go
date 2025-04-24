package ezapp

import "context"

// AppOption is a functional option for configuring an App.
type AppOption func(*App)

// WithRunnables adds one or more runnables to the App.
func WithRunnables(runnables ...Runnable) AppOption {
	return func(app *App) {
		app.runnables = append(app.runnables, runnables...)
	}
}

// WithCleanup sets the cleanup function for the App.
func WithCleanup(cleanup func(context.Context) error) AppOption {
	return func(app *App) {
		app.cleanup = cleanup
	}
}

// Construct creates a new App with the provided options.
func Construct(options ...AppOption) App {
	app := App{}
	for _, option := range options {
		option(&app)
	}
	return app
}
