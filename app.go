package copper

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tusharsoni/copper/cconfig"
	"github.com/tusharsoni/copper/clogger"
)

// Runner provides an interface that can be run by a Copper app using the Run or Start funcs.
// This interface is implemented by various packages within Copper including chttp.Server.
type Runner interface {
	Run() error
}

// NewApp creates a new Copper app and returns it along with the app's lifecycle manager,
// config, and the logger.
func NewApp(lifecycle *Lifecycle, config cconfig.Config, logger clogger.Logger) *App {
	return &App{
		Lifecycle: lifecycle,
		Config:    config,
		Logger:    logger,
	}
}

// App defines a Copper app container that can run provided code in its managed lifecycle.
// It provides functionality to read config in multiple environments as defined by
// command-line flags.
type App struct {
	Lifecycle *Lifecycle
	Config    cconfig.Config
	Logger    clogger.Logger
}

// Run runs the provided funcs. Once all of the functions complete their run,
// the  lifecycle's stop funcs are also called. If any of the fns return an error,
// the app exits with an exit code 1.
// Run should be used when none of the fn are blocking. For blocking funcs like
// a long running server, use Start.
func (a *App) Run(fns ...Runner) {
	for i := range fns {
		err := fns[i].Run()
		if err != nil {
			a.Logger.Error("Failed to run", err)
			a.Lifecycle.Stop()
			os.Exit(1)
		}
	}

	a.Lifecycle.Stop()
}

// Start runs the provided fns where the last provided fn be blocking like a
// long running  server. The app listens on the OS's INT and TERM signals from the
// user to exit. Once the signal is received, the lifecycle's stop funcs are
// called.
// If any of the fns fail to run and returns an error, the app exits with exit code
// 1.
// If none of the fns are blocking, use Run instead.
func (a *App) Start(fns ...Runner) {
	osInt := make(chan os.Signal, 1)

	signal.Notify(osInt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-osInt
		a.Lifecycle.Stop()
	}()

	for i := range fns {
		err := fns[i].Run()
		if err != nil {
			a.Logger.Error("Failed to run", err)
			os.Exit(1)
		}
	}
}
