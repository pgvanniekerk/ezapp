package wire

type Impl[IFace any] struct{}

func (i *Impl[IFace]) Implements() IFace {
	return *new(IFace)
}
