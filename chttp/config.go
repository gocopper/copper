package chttp

import (
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cerrors"
)

// LoadConfig loads Config from app's config
func LoadConfig(appConfig cconfig.Loader) (Config, error) {
	var config Config

	err := appConfig.Load("chttp", &config)
	if err != nil {
		return Config{}, cerrors.New(err, "failed to load chttp config", nil)
	}

	return config, nil
}

// Config holds the params needed to configure Server
type Config struct {
	Port                    uint `default:"7501"`
	DevMode                 bool `toml:"dev_mode"`
	EnableSinglePageRouting bool `toml:"enable_single_page_routing"`
}
