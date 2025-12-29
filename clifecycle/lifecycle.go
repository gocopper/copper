package clifecycle

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gocopper/copper/clogger"
)

const defaultStopTimeout = 30 * time.Second

// New instantiates and returns a new Lifecycle that can be used with
// New to create a Copper app.
func New(logger clogger.Logger) *Lifecycle {
	ctx, cancel := context.WithCancel(context.Background())

	return &Lifecycle{
		Context:     ctx,
		logger:      logger,
		cancel:      cancel,
		onStop:      make([]func(ctx context.Context) error, 0),
		stopTimeout: defaultStopTimeout,
	}
}

// Lifecycle represents the lifecycle of an app. Most importantly, it
// allows various parts of the app to register stop funcs that are run
// before the app exits.
// Packages such as chttp use Lifecycle to gracefully stop the HTTP
// server before the app exits.
type Lifecycle struct {
	Context     context.Context
	logger      clogger.Logger
	cancel      context.CancelFunc
	onStop      []func(ctx context.Context) error
	stopTimeout time.Duration
	wg          sync.WaitGroup
}

// OnStop registers the provided fn to run before the app exits. The fn
// is given a context with a deadline. Once the deadline expires, the
// app may exit forcefully.
func (lc *Lifecycle) OnStop(fn func(ctx context.Context) error) {
	lc.onStop = append(lc.onStop, fn)
}

// Go starts a background goroutine that will be waited for during shutdown.
// The goroutine should return when the context is done or when its work is complete.
//
// WARNING: Be careful with closure capture in loops. Make copies of loop variables:
//
//	for _, handler := range handlers {
//	    handler := handler  // copy to avoid closure capture bug
//	    lc.Go(func(ctx context.Context) {
//	        err := handler.Process(ctx, payload)
//	    })
//	}
func (lc *Lifecycle) Go(fn func(ctx context.Context)) {
	lc.wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				lc.logger.WithTags(map[string]interface{}{
					"error": r,
					"stack": string(debug.Stack()),
				}).Error("[copper] Panic in goroutine", nil)
			}
		}()
		defer lc.wg.Done()
		fn(lc.Context)
	}()
}

// Stop runs all of the registered stop funcs in order along with a
// context with a configured timeout and waits for them to complete.
func (lc *Lifecycle) Stop(logger Logger) {
	// Cancel context first so goroutines know to stop
	lc.cancel()

	// Create shutdown context with timeout for cleanup operations
	shutdownCtx, cancel := context.WithTimeout(context.Background(), lc.stopTimeout)
	defer cancel()

	// Run cleanup functions (HTTP server shutdown, etc.)
	for _, fn := range lc.onStop {
		err := fn(shutdownCtx)
		if err != nil {
			logger.Error("Failed to run cleanup func", err)
		}
	}

	// Wait for background goroutines to complete with timeout
	done := make(chan struct{})
	go func() {
		lc.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("All background jobs completed successfully")
	case <-shutdownCtx.Done():
		logger.Error("Background jobs did not complete within timeout", shutdownCtx.Err())
	}
}
