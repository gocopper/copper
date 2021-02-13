package copper

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/tusharsoni/copper/v2/cconfig"
	"github.com/tusharsoni/copper/v2/clogger"
)

// Start starts the app by calling the provided function with an app context, logger,
// and config. It registers command line flags to set config path and environment.
// The app runs until it receives an interrupt signal from the OS.
func Start(f func(ctx context.Context, logger clogger.Logger, config cconfig.Config) error) {
	var (
		env       = flag.String("env", "dev", "Current environment (dev, test, staging, prod)")
		configDir = flag.String("config", "./config", "Path to directory with config files")
	)

	flag.Parse()

	config, err := cconfig.NewConfig(*configDir, *env)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := clogger.NewConsoleWithConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	osInt := make(chan os.Signal, 1)

	signal.Notify(osInt, os.Interrupt)

	go func() {
		<-osInt
		cancel()
	}()

	err = f(ctx, logger, config)
	if err != nil {
		log.Fatal(err)
	}
}
