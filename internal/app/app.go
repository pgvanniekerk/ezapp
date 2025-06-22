package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func New(runnerList []Runner, logger *zap.Logger) App {
	return App{
		runnerList: runnerList,
		logger:     logger,
	}
}

type App struct {
	runnerList []Runner
	logger     *zap.Logger
}

func (a App) Run() error {
	a.logger.Debug("start application")

	// Create a termination context with a cancel function that is
	// used to signal application termination.
	termCtx, termFunc := context.WithCancel(context.Background())
	defer termFunc()
	a.logger.Debug("created termination context")

	// Asynchronously listen for SIGINT, SIGTERM. If signaled,
	// the termCtx will be canceled and propagated to all runnable
	// invocations.
	go a.terminationSignaller(termFunc)
	a.logger.Debug("started termination signaller")

	// Create an error group with context that will be used to
	// asynchronously invoke each runnable.
	// Should an error occur, the error group will automatically
	// cancel the context, propagating to each runnable - starting
	// the shutdown process.
	errGrp, ctx := errgroup.WithContext(termCtx)
	a.logger.Debug("created error group")

	// Invoke each runnable through the error group.
	for idx, _ := range a.runnerList {
		errGrp.Go(func() error {
			return a.runnerList[idx](ctx)
		})
	}
	a.logger.Debug("started runnable invocations via error group")

	// Wait for an error or for all runnable invocations to finalize
	// and return.
	err := errGrp.Wait()
	if err != nil {
		return fmt.Errorf("failed to invoke runnable: %w", err)
	}
	a.logger.Debug("application finished running")

	return nil
}

// terminationSignaller is a helper function that listens for SIGINT and SIGTERM
// and cancels the given termFunc.
func (a App) terminationSignaller(termFunc context.CancelFunc) {
	a.logger.Debug("starting termination signaller")

	// Listen for SIGINT and SIGTERM and notify via sigChan.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	a.logger.Debug("started listening for SIGINT and SIGTERM")

	// Wait for signal then cancel termCtx.
	<-sigChan
	termFunc()
	a.logger.Debug("received SIGINT or SIGTERM, terminating")

	// Free/Release signal processing objects.
	signal.Stop(sigChan)
	close(sigChan)
	a.logger.Debug("stopped listening for SIGINT and SIGTERM")

}
