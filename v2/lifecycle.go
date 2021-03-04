package copper

import (
	"context"
	"time"

	"github.com/tusharsoni/copper/v2/clogger"
)

const defaultStopTimeout = 10 * time.Second

// NewLifecycle instantiates and returns a new Lifecycle that can be used with
// New to create a Copper app.
func NewLifecycle(logger clogger.Logger) *Lifecycle {
	return &Lifecycle{
		logger:      logger,
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
	logger      clogger.Logger
	onStop      []func(ctx context.Context) error
	stopTimeout time.Duration
}

// OnStop registers the provided fn to run before the app exits. The fn
// is given a context with a deadline. Once the deadline expires, the
// app may exit forcefully.
func (lc *Lifecycle) OnStop(fn func(ctx context.Context) error) {
	lc.onStop = append(lc.onStop, fn)
}

// Stop runs all of the registered stop funcs in order along with a
// context with a configured timeout.
func (lc *Lifecycle) Stop() {
	for _, fn := range lc.onStop {
		ctx, cancel := context.WithTimeout(context.Background(), lc.stopTimeout)

		err := fn(ctx)
		if err != nil {
			lc.logger.Error("Failed to run cleanup func", err)
		}

		cancel()
	}
}
