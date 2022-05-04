package copper

import (
	"flag"

	"github.com/gocopper/copper/cconfig"
)

// Flags holds flag values passed in via command line. These can be used to configure the app environment
// and override the config directory.
type Flags struct {
	ConfigPath cconfig.Path
}

// NewFlags reads the command line flags and returns Flags with the values set.
func NewFlags() *Flags {
	var configPath = flag.String("config", "./config/local.toml", "Path to config file")

	flag.Parse()

	return &Flags{
		ConfigPath: cconfig.Path(*configPath),
	}
}
