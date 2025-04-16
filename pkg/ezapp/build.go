package ezapp

import (
	"errors"
	"fmt"
	"reflect"

	env "github.com/Netflix/go-env"
)

func Build[CONF any](builder Builder[CONF]) EzApp {

	var conf CONF

	// Validate that CONF is a struct
	if reflect.TypeOf(conf).Kind() != reflect.Struct {
		return EzApp{
			initErr: errors.New("CONF must be a struct"),
		}
	}

	// Use go-env to populate CONF from environment variables
	if _, err := env.UnmarshalFromEnviron(&conf); err != nil {
		return EzApp{
			initErr: fmt.Errorf("failed to parse environment variables into CONF: %w", err),
		}
	}

	// Call the builder function to get the list of runnables
	runnables, err := builder(conf)
	if err != nil {
		return EzApp{
			initErr: err,
		}
	}

	// Create and return a new EzApp with the runnables
	return EzApp{
		runnableList: runnables,
	}
}
