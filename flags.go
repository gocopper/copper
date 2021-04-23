package copper

import (
	"flag"

	"github.com/tusharsoni/copper/cconfig"
)

// Flags holds flag values passed in via command line. These can be used to configure the app environment
// and override the config directory.
type Flags struct {
	Env        cconfig.Env
	ConfigDir  cconfig.Dir
	ProjectDir cconfig.ProjectDir
}

// NewFlags reads the command line flags and returns Flags with the values set.
func NewFlags() *Flags {
	var (
		env        = flag.String("env", "dev", "Current environment (dev, test, staging, prod)")
		configDir  = flag.String("config", "./config", "Path to directory with config files")
		projectDir = flag.String("project", ".", "Path to project directory")
	)

	flag.Parse()

	return &Flags{
		Env:        cconfig.Env(*env),
		ConfigDir:  cconfig.Dir(*configDir),
		ProjectDir: cconfig.ProjectDir(*projectDir),
	}
}
