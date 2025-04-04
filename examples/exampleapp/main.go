package main

import (
	"github.com/pgvanniekerk/ezapp/pkg/buildoption"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"log"
	"time"
)

// CustomErrorHandler logs errors instead of panicking
func CustomErrorHandler(err error) error {
	log.Printf("Error occurred: %v", err)
	return err
}

func main() {
	ezapp.Build(
		WireFunc,
		buildoption.WithOptions(
			buildoption.WithErrorHandler(CustomErrorHandler),
			buildoption.WithStartupTimeout(30*time.Second),
			buildoption.WithEnvVarPrefix("EXAMPLEAPP"),
		),
	).Run()
}
