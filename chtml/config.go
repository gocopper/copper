package chtml

import (
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

// LoadConfig loads Config from app config
func LoadConfig(appConfig cconfig.Loader) (Config, error) {
	var config Config

	err := appConfig.Load("chtml", &config)
	if err != nil {
		return Config{}, cerrors.New(err, "failed to load chtml config", nil)
	}

	return config, nil
}

// Config holds params to configure chtml
type Config struct {
	DevMode                 bool `toml:"dev_mode"`
	EnableSinglePageRouting bool `toml:"enable_single_page_routing"`
}
