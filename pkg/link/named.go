package link

import "go.uber.org/dig"

// Named returns an Option that specifies the name of a constructor.
func Named(name string) Option {
	return func(opts *[]dig.ProvideOption) {
		*opts = append(*opts, dig.Name(name))
	}
}
