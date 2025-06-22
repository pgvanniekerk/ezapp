# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

EzApp is a lightweight Go framework for building applications with:
- Environment-based configuration loading
- Concurrent service execution with graceful shutdown
- Dependency injection through builder functions
- Signal handling (SIGINT/SIGTERM) with context propagation

## Development Commands

### Build and Test
```bash
go build ./...          # Build all packages
go test ./...           # Run all tests
go test -v ./...        # Run tests with verbose output
go test ./internal/config  # Run tests for specific package
```

### Code Quality
```bash
go fmt ./...            # Format code
go vet ./...            # Static analysis
go mod tidy             # Clean up dependencies
```

### Running
```bash
go run main.go          # Run main package (if exists)
go run ./cmd/myapp      # Run specific command
```

## Architecture

### Core Components

**ezapp.go** - Main entry point with generic `Run[Config]` function that:
- Takes an `Initializer[Config]` function for dependency injection
- Returns `AppCtx` containing list of `Runner` functions
- Provides `InitCtx[Config]` with startup context, logger, and config

**internal/app/app.go** - Application orchestration that:
- Manages concurrent execution of multiple `Runner` functions
- Handles graceful shutdown via context cancellation
- Uses `golang.org/x/sync/errgroup` for coordinated error handling
- Listens for SIGINT/SIGTERM signals

**internal/app/runner.go** - Defines `Runner` interface:
```go
type Runner func(context.Context) error
```

**internal/config/** - Configuration utilities for:
- Loading environment variables with defaults
- Creating startup contexts with timeouts
- Logger initialization

### Current Refactoring State

The codebase is undergoing a major refactoring (see git status). The old dependency injection system using:
- `wire` package components
- `container` package for DI
- `link` package for component linking

Is being replaced with a simpler system using:
- Generic `Initializer[Config]` functions
- `Runner` function types instead of `Runnable` interface
- Direct function composition for dependency injection

### Key Patterns

1. **Configuration**: Use structs with `env` tags for environment variable binding
2. **Services**: Implement as functions that return `Runner` functions
3. **Dependency Injection**: Use the `Initializer` function to wire dependencies
4. **Graceful Shutdown**: Services should respect context cancellation
5. **Error Handling**: Return errors from `Runner` functions for application-level issues

### Testing

Tests are located in `internal/config/*_test.go` and follow standard Go testing conventions. Tests use table-driven patterns and environment variable manipulation for configuration testing.

### Dependencies

- `go.uber.org/zap` - Structured logging
- `github.com/Netflix/go-env` - Environment variable loading
- `golang.org/x/sync/errgroup` - Coordinated goroutine execution
- `github.com/stretchr/testify` - Testing utilities