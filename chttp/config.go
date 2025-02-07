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
	Port                               uint    `default:"7501"`
	UseLocalHTML                       bool    `toml:"use_local_html"`
	RenderHTMLError                    bool    `toml:"render_html_error"`
	EnableSinglePageRouting            bool    `toml:"enable_single_page_routing"`
	ReadTimeoutSeconds                 uint    `toml:"read_timeout_seconds" default:"10"`
	RedirectURLForUnauthorizedRequests *string `toml:"redirect_url_for_unauthorized_requests"`
	BasePath                           *string `toml:"base_path"`
}
