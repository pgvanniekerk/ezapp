package wire

import "context"

type Component struct{}

func (c *Component) Init(_ context.Context, _ struct{}) error {
	return nil
}

func (c *Component) Cleanup(_ context.Context) error {
	return nil
}
