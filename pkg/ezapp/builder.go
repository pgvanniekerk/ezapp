package ezapp

type Builder[CONF any] func(CONF) ([]Runnable, error)
