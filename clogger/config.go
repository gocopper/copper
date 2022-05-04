package clogger

import (
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

// Format represents the output format of log statements
type Format string

// Formats supported by Logger
const (
	FormatPlain = Format("plain")
	FormatJSON  = Format("json")
)

// LoadConfig loads Config from app's config
func LoadConfig(appConfig cconfig.Loader) (Config, error) {
	var config Config

	err := appConfig.Load("clogger", &config)
	if err != nil {
		return Config{}, cerrors.New(err, "failed to load clogger config", nil)
	}

	if config.Format != FormatPlain && config.Format != FormatJSON {
		config.Format = FormatPlain
	}

	return config, nil
}

// Config holds the params needed to configure Logger
type Config struct {
	Out    string `toml:"out"`
	Err    string `toml:"err"`
	Format Format `toml:"format"`
}
