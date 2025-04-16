package ezapp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// osExit is a package-level variable that can be overridden during tests
var osExit = os.Exit

type EzApp struct {
	runnableList []Runnable
	initErr      error
}

func (e EzApp) Run() {

	// Check if there was an initialization error
	if e.initErr != nil {
		fmt.Printf("App shutting down: Initialization error: %v\n", e.initErr)
		osExit(1)
	}

	// Check if there are any runnables to run
	if len(e.runnableList) == 0 {
		fmt.Println("App shutting down: No runnables to execute")
		osExit(0)
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan) // Stop signal handling when done

	// Create error channel to collect errors from runnables
	errChan := make(chan error, len(e.runnableList))

	// Create a wait group to wait for all runnables to finish
	var wg sync.WaitGroup

	// Start all runnables
	for _, runnable := range e.runnableList {
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
	}

	// Wait for all runnables to finish
	<-done

	// Close the error channel
	close(errChan)

	fmt.Println("App shutdown complete")
	osExit(exitCode)
}
