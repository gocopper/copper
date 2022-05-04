package copper

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/clogger"
)

// Runner provides an interface that can be run by a Copper app using the Run or Start funcs.
// This interface is implemented by various packages within Copper including chttp.Server.
type Runner interface {
	Run() error
}

// New provides a convenience wrapper around InitApp that logs and exits if there is an error.
func New() *App {
	app, err := InitApp()
	if err != nil {
		log.Fatal(err)
	}

	return app
}

// NewApp creates a new Copper app and returns it along with the app's lifecycle manager,
// config, and the logger.
func NewApp(lifecycle *Lifecycle, config cconfig.Loader, logger clogger.Logger) *App {
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
	Config    cconfig.Loader
	Logger    clogger.Logger
}

// Run runs the provided funcs. Once all of the functions complete their run,
// the  lifecycle's stop funcs are also called. If any of the fns return an error,
// the app exits with an exit code 1.
// Run should be used when none of the fn are long-running. For long-running funcs like
// an HTTP server, use Start.
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

// Start runs the provided fns and then waits on the OS's INT and TERM signals from the
// user to exit. Once the signal is received, the lifecycle's stop funcs are
// called.
// If any of the fns fail to run and returns an error, the app exits with exit code
// 1.
func (a *App) Start(fns ...Runner) {
	for i := range fns {
		err := fns[i].Run()
		if err != nil {
			a.Logger.Error("Failed to run", err)
			os.Exit(1)
		}
	}

	osInt := make(chan os.Signal, 1)

	signal.Notify(osInt, syscall.SIGINT, syscall.SIGTERM)

	<-osInt

	a.Lifecycle.Stop()
}
