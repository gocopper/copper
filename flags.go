package copper

import (
	"flag"

	"github.com/gocopper/copper/cconfig"
)

// Flags holds flag values passed in via command line. These can be used to configure the app environment
// and override the config directory.
type Flags struct {
	ConfigPath      cconfig.Path
	ConfigOverrides cconfig.Overrides
}

// NewFlags reads the command line flags and returns Flags with the values set.
func NewFlags() *Flags {
	var (
		configPath      = flag.String("config", "./config/dev.toml", "Path to config file")
		configOverrides = flag.String("set", "", "Config overrides ex. \"chttp.port=5902\". Separate multiple overrides with ;")
	)

	flag.Parse()

	return &Flags{
		ConfigPath:      cconfig.Path(*configPath),
		ConfigOverrides: cconfig.Overrides(*configOverrides),
	}
}
