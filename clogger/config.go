package clogger

import "os"

// StdConfig can be used to customize the setup for the std logger.
type StdConfig struct {
	MinLevel Level
}

func (c StdConfig) isValid() bool {
	return c.MinLevel > 0
}

// GetDefaultStdConfig provides a config for the std logger with defaults
func GetDefaultStdConfig() StdConfig {
	return StdConfig{
		MinLevel: LevelDebug,
	}
}

// SentryConfig holds the config values for the sentry logger
type SentryConfig struct {
	Dsn                   string
	MinLevelForStd        Level
	MinLevelForCapture    Level
	IgnoredErrsForCapture []error
}

func (c SentryConfig) isValid() bool {
	return len(c.Dsn) > 0 &&
		c.MinLevelForStd > 0 &&
		c.MinLevelForCapture > 0
}

// GetDefaultSentryConfig provides a default set of config values for sentry logger
func GetDefaultSentryConfig() SentryConfig {
	return SentryConfig{
		Dsn:                   os.Getenv("Sentry_DSN"),
		MinLevelForStd:        LevelDebug,
		MinLevelForCapture:    LevelWarn,
		IgnoredErrsForCapture: []error{},
	}
}
