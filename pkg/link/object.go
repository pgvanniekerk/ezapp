package link

import (
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"go.uber.org/dig"
)

func Object[T any](obj T, opt Option) ezapp.BuildProcess {

	opts := make([]dig.ProvideOption, 0)

	if opt != nil {
		opt(&opts)
	}

	return func(bCtx ezapp.BuildContext) error {
		return bCtx.Container().Provide(func() T {
			return obj
		}, opts...)
	}
}
