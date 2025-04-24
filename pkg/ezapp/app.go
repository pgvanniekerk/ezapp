package ezapp

import "context"

// App represents an application with runnables and cleanup functions.
type App struct {
	runnables []Runnable
	cleanup   func(context.Context) error
}
