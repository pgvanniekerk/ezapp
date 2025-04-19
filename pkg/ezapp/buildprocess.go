package ezapp

// BuildProcess is a functional option for configuring the application
type BuildProcess func(BuildContext) error
