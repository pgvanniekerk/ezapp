package ezapp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/Netflix/go-env"
)

var osExit = os.Exit

// Run creates an application from the provided builder function and runs it.
// It executes the builder to get runnables, then starts all runnables while listening for stop signals.
func Run[CONF any](builder Builder[CONF]) {
	var conf CONF

	// Validate that CONF is a struct
	if reflect.TypeOf(conf).Kind() != reflect.Struct {
		fmt.Printf("App shutting down: Initialization error: %v\n", errors.New("CONF must be a struct"))
		osExit(1)
		return
	}

	// Use go-env to populate CONF from environment variables
	if _, err := env.UnmarshalFromEnviron(&conf); err != nil {
		fmt.Printf("App shutting down: Initialization error: %v\n", fmt.Errorf("failed to parse environment variables into CONF: %w", err))
		osExit(1)
		return
	}

	// Call the builder function to get the list of runnables
	runnables, err := builder(conf)
	if err != nil {
		fmt.Printf("App shutting down: Initialization error: %v\n", err)
		osExit(1)
		return
	}

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
	defer signal.Stop(sigChan) // Stop signal handling when done

	// Create error channel to collect errors from runnables
	errChan := make(chan error, len(runnables))

	// Create a wait group to wait for all runnables to finish
	var wg sync.WaitGroup

	// Start all runnables
	for _, runnable := range runnables {
		wg.Add(1)
		go func(r Runnable) {
			defer wg.Done()
			if err := r.Run(ctx); err != nil {
				select {
				case errChan <- err:
					// Error sent successfully
				case <-ctx.Done():
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

	// Wait for either an error, SIGINT, SIGTERM, or all runnables to finish
	var exitCode int
	select {
	case sig := <-sigChan:
		// Received SIGINT or SIGTERM, cancel context
		fmt.Printf("App shutting down: Received signal %s\n", sig.String())
		cancel()
		exitCode = 0 // Normal termination due to signal
	case err := <-errChan:
		// Received an error from a runnable
		fmt.Printf("App shutting down: Runnable error: %v\n", err)
		// Context is already cancelled in the goroutine
		exitCode = 1 // Error termination
	case <-done:
		// All runnables finished successfully
		fmt.Println("App shutting down: All runnables completed successfully")
		osExit(0) // Exit immediately with success code
		return
	}

	// Wait for all runnables to finish
	<-done

	// Close the error channel
	close(errChan)

	fmt.Println("App shutdown complete")
	osExit(exitCode)
}
