package primitive

import "context"

type Server[Params any] interface {
	Component[Params]
	Start(context.Context) error
	Stop(context.Context) error
}
