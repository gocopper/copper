package chttp

type config struct {
	Port                uint `default:"7501"`
	ShutdownTimeoutSecs uint `default:"15"`
}
