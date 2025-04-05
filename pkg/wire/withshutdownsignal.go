package wire

// WithShutdownSignal returns an AppOption that sets the shutdown signal channel for the App
func WithShutdownSignal(shutdownSig <-chan error) AppOption {
	return func(o *appOptions) {
		o.shutdownSig = shutdownSig
	}
}
