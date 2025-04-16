package ezapp

import "context"

type Runnable interface {
	Run(context.Context) error
}
