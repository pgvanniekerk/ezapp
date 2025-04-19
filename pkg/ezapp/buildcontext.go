package ezapp

import (
	"context"
	"go.uber.org/dig"
)

type BuildContext interface {
	Container() *dig.Container
	InitTimeout() context.Context
	Modules() []Module
}

// buildContext implements the BuildContext interface
type buildContext struct {
	container *dig.Container
	initCtx   context.Context
}

// Container returns the dig container
func (b *buildContext) Container() *dig.Container {
	return b.container
}

func (b *buildContext) InitTimeout() context.Context {
	return b.initCtx
}
