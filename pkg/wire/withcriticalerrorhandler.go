package wire

// WithCriticalErrHandler sets the handler for critical errors from runnables.
// This handler will be called when a runnable reports a critical error.
//
// Example:
//
//	app, err := wire.App(
//	    wire.Runnables(myRunnable),
//	    wire.WithCriticalErrHandler(func(err error) {
//	        log.Fatalf("Critical error: %v", err)
//	    }),
//	)
func WithCriticalErrHandler(handler func(error)) AppOption {
	return func(opts *appOptions) {
		opts.criticalErrHandler = handler
	}
}
