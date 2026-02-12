package clifecycle

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gocopper/copper/clogger"
)

const defaultStopTimeout = 30 * time.Second

func New(logger clogger.CoreLogger) *Lifecycle {
	ctx, cancel := context.WithCancel(context.Background())

	return &Lifecycle{
		ctx:         ctx,
		logger:      logger,
		cancel:      cancel,
		onStop:      make([]func(ctx context.Context) error, 0),
		stopTimeout: defaultStopTimeout,
		shutdown:    make(chan struct{}),
	}
}

type Lifecycle struct {
	ctx         context.Context
	logger      clogger.CoreLogger
	cancel      context.CancelFunc
	onStop      []func(ctx context.Context) error
	stopTimeout time.Duration
	wg          sync.WaitGroup
	shutdown    chan struct{}
}

// Context returns the lifecycle context which remains valid until shutdown completes.
func (lc *Lifecycle) Context() context.Context {
	return lc.ctx
}

// Shutdown returns a channel that's closed when graceful shutdown begins.
// Goroutines should select on this channel to exit their loops gracefully.
func (lc *Lifecycle) Shutdown() <-chan struct{} {
	return lc.shutdown
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
		fn(lc.ctx)
	}()
}

// Stop signals graceful shutdown and waits for background goroutines to complete.
// Context remains valid during graceful shutdown. Only canceled if timeout expires (forced shutdown).
func (lc *Lifecycle) Stop(logger Logger) {
	// Signal graceful shutdown
	close(lc.shutdown)

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), lc.stopTimeout)
	defer cancel()

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

	// Run cleanup functions (close DB, etc.)
	for _, fn := range lc.onStop {
		err := fn(shutdownCtx)
		if err != nil {
			logger.Error("Failed to run cleanup func", err)
		}
	}

	// Cancel lifecycle context for final cleanup
	lc.cancel()
}
