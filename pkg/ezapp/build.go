package ezapp

import (
	"io"
	"log/slog"

	"github.com/pgvanniekerk/ezapp/internal/app"
)

func Build(optF optionsFunc, runnables ...Runnable) EzApp {

	opts := new(options)
	for _, optFunc := range optF() {
		optFunc(opts)
	}

	// Create a noop logger if no logger has been passed
	if opts.logger == nil {
		opts.logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	return app.New(
		runnables,
		opts.logger,
	)
}
