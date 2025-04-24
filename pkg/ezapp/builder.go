package ezapp

type Builder[CONF any] func(CONF) (App, error)
