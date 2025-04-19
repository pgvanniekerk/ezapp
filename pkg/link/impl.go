package link

import "go.uber.org/dig"

// Impl returns an Option that specifies that a constructor provides an implementation of the given interface.
func Impl[T any](opt Option) Option {
	return func(opts *[]dig.ProvideOption) {
		opt(opts)
		*opts = append(*opts, dig.As(new(T)))
	}
}
