package app

import "context"

type Runnable interface {
	Run(context.Context) error
	HandleError(error) error
}
