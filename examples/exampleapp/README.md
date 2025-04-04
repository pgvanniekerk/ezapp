# Example App

This example demonstrates a complete application using ezapp with all available options. It includes:

- A compact `main()` function
- Mock database connection using the startup context
- A mock service that uses the mock database
- Cleanup function for the mock database
- Custom error handler
- Custom startup timeout
- Custom environment variable prefix

## Overview

The example consists of three files:

1. `main.go`: Contains the compact main function that builds and runs the application
2. `wire.go`: Contains the configuration struct and wire function for dependency injection
3. `service.go`: Contains the mock service implementation

The application creates a mock database, connects to it using the startup context, creates a mock service that periodically queries the mock database, and sets up proper cleanup for the mock database connection when the application shuts down.

## Running the Example

Running the example is simple:

```bash
# Run the example
go run .
```

The example uses a mock database, so no external database is required.

## Key Points

### Compact Main Function

The main function is kept compact and focused on building and running the application:

```go
func main() {
    ezapp.Build(
        WireFunc,
        buildoption.WithOptions(
            buildoption.WithErrorHandler(CustomErrorHandler),
            buildoption.WithStartupTimeout(30 * time.Second),
            buildoption.WithEnvVarPrefix("EXAMPLEAPP"),
        ),
    ).Run()
}
```

### Dependency Injection with WireFunc

The `WireFunc` function demonstrates proper dependency injection:

1. It creates a mock database and connects to it using the startup context
2. It creates a mock service that depends on the mock database
3. It provides a cleanup function to close the mock database connection

### Service Implementation

The `MockService` implements the `Service` interface with:

1. A `Run()` method that starts the service and returns an error only in exceptional circumstances such as dependency failures or timeouts (application-impacting errors)
2. A `Stop(context.Context)` method that gracefully shuts down the service, respecting the context timeout. If it returns an error, it will be reported during shutdown. If the context timeout is reached, the application will force close

### Error Handling

The example includes a custom error handler that logs errors instead of panicking:

```go
func CustomErrorHandler(err error) error {
    log.Printf("Error occurred: %v", err)
    return err
}
```

## What This Example Demonstrates

- How to structure a complete application using ezapp
- How to use all available options in a compact main function
- How to properly wire dependencies using the startup context
- How to implement the Service interface
- How to handle cleanup of resources
