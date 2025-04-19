package link

import "go.uber.org/dig"

// Option is a function that modifies a slice of dig.ProvideOption
type Option func(*[]dig.ProvideOption)
