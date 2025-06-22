package app

import "context"

type Runner func(context.Context) error
