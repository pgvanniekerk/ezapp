package ezapp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Netflix/go-env"
)

var osExit = os.Exit

// GetCleanupTimeout returns the timeout duration for the cleanup function.
// It retrieves the timeout in seconds from the EZAPP_TERM_TIMEOUT environment variable.
// If not found, it defaults to 15 seconds.
func GetCleanupTimeout() time.Duration {
	timeoutStr := os.Getenv("EZAPP_TERM_TIMEOUT")
	if timeoutStr == "" {
		return 15 * time.Second
	}

	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil {
		// If the value is not a valid integer, use the default
		return 15 * time.Second
	}

	return time.Duration(timeoutSec) * time.Second
}

// validateAndLoadConfig validates that CONF is a struct and loads environment variables into it.
func validateAndLoadConfig[CONF any](conf *CONF) error {
	// Validate that CONF is a struct
	if reflect.TypeOf(*conf).Kind() != reflect.Struct {
		return errors.New("CONF must be a struct")
	}

	// Use go-env to populate CONF from environment variables
	if _, err := env.UnmarshalFromEnviron(conf); err != nil {
		return fmt.Errorf("failed to parse environment variables into CONF: %w", err)
	}

	return nil
}

// startRunnables starts all runnables in separate goroutines and returns channels for errors and completion.
func startRunnables(ctx context.Context, runnables []Runnable) (chan error, chan struct{}, context.CancelFunc) {
	// Create a cancellable context
	runCtx, cancel := context.WithCancel(ctx)

	// Create error channel to collect errors from runnables
	errChan := make(chan error, len(runnables))

	// Create a wait group to wait for all runnables to finish
	var wg sync.WaitGroup

	// Start all runnables
	for _, runnable := range runnables {
		wg.Add(1)
		go func(r Runnable) {
			defer wg.Done()
			if err := r.Run(runCtx); err != nil {
				select {
				case errChan <- err:
					// Error sent successfully
				case <-runCtx.Done():
					// Context already cancelled, no need to send error
				}
				cancel() // Cancel context on error
			}
		}(runnable)
	}

	// Create a channel to signal when all runnables are done
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	return errChan, done, cancel
}

// waitForShutdown waits for either a signal, an error, or all runnables to finish.
func waitForShutdown(sigChan chan os.Signal, errChan chan error, done chan struct{}, cancel context.CancelFunc) (int, string) {
	// Wait for either an error, SIGINT, SIGTERM, or all runnables to finish
	var exitCode int
	var shutdownReason string
	select {
	case sig := <-sigChan:
		// Received SIGINT or SIGTERM, cancel context
		shutdownReason = fmt.Sprintf("Received signal %s", sig.String())
		cancel()
		exitCode = 0 // Normal termination due to signal
	case err := <-errChan:
		// Received an error from a runnable
		shutdownReason = fmt.Sprintf("Runnable error: %v", err)
		// Context is already cancelled in the goroutine
		exitCode = 1 // Error termination
	case <-done:
		// All runnables finished successfully
		shutdownReason = "All runnables completed successfully"
		exitCode = 0 // Success
	}

	return exitCode, shutdownReason
}

// runCleanup runs the cleanup function with a timeout.
func runCleanup(cleanup func(context.Context) error) (int, error) {
	fmt.Println("Running cleanup function")
	// Create a context with timeout for the cleanup function
	timeout := GetCleanupTimeout()
	cleanupCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := cleanup(cleanupCtx)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

// Run creates an application from the provided builder function and runs it.
// It executes the builder to get an App, then starts all runnables while listening for stop signals.
func Run[CONF any](builder Builder[CONF]) {
	var conf CONF

	// Validate and load configuration
	if err := validateAndLoadConfig(&conf); err != nil {
		fmt.Printf("App shutting down: Initialization error: %v\n", err)
		osExit(1)
		return
	}

	// Call the builder function to get the App
	app, err := builder(conf)
	if err != nil {
		fmt.Printf("App shutting down: Initialization error: %v\n", err)
		osExit(1)
		return
	}

	// Get the runnables from the App
	runnables := app.runnables

	// Check if there are any runnables to run
	if len(runnables) == 0 {
		fmt.Println("App shutting down: No runnables to execute")
		osExit(0)
		return
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start all runnables
	errChan, done, runnablesCancel := startRunnables(ctx, runnables)

	// Wait for shutdown
	exitCode, shutdownReason := waitForShutdown(sigChan, errChan, done, runnablesCancel)

	fmt.Printf("App shutting down: %s\n", shutdownReason)

	// If we're not already done, wait for all runnables to finish
	if shutdownReason != "All runnables completed successfully" {
		<-done
	}

	// Close the error channel
	close(errChan)

	// Call the cleanup function if it exists
	if app.cleanup != nil {
		cleanupExitCode, err := runCleanup(app.cleanup)
		if err != nil {
			fmt.Printf("Cleanup error: %v\n", err)
			exitCode = cleanupExitCode
		}
	}
	signal.Stop(sigChan) // Stop signal handling when done
	close(sigChan)
	runnablesCancel()

	fmt.Println("App shutdown complete")
	osExit(exitCode)
}
