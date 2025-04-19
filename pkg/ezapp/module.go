package ezapp

import "reflect"

type Module interface {
	LoadModule(BuildContext) error
	GetDependencies(reflect.Type) (reflect.Type, error)
	WireDependencies(reflect.Value, reflect.Value) error
}
