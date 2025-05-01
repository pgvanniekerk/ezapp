package primitive

import "context"

type Service[Params any] interface {
	Component[Params]
	Start(context.Context) error
	Stop(context.Context) error
}
