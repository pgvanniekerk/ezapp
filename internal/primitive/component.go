package primitive

import (
	"context"
	"go.uber.org/dig"
)

type Component[Params any] interface {
	Init(context.Context, Params) error
	Cleanup(context.Context) error
	Wire(*dig.Container) error
}
