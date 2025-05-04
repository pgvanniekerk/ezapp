package primitive

import (
	"context"
)

type Component[Params any] interface {
	Init(context.Context, Params) error
	Cleanup(context.Context) error
}
