package link

import "go.uber.org/dig"

// Grouped returns an Option that specifies the group of a constructor.
func Grouped(group string) Option {
	return func(opts *[]dig.ProvideOption) {
		*opts = append(*opts, dig.Group(group))
	}
}
