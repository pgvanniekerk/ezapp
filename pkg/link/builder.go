package link

import (
	"context"
	"fmt"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

func Builder[B builder[T], T any](bCtx ezapp.BuildContext) error {

	err := bCtx.Container().Provide(
		func(b B) (T, error) {
			t, err := b.Build(bCtx.InitTimeout())
			return t, fmt.Errorf("error building %T: %w", t, err)
		},
	)
	if err != nil {
		return err
	}

	return nil
}

type builder[T any] interface {
	Build(context.Context) (T, error)
}
