package ezapp

type optionsFunc func() []func(*options)

func WithOptions(funcs ...func(*options)) []func(*options) {
	return funcs
}
