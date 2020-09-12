package cqueue

import "time"

type Config struct {
	DequeueWaitTime time.Duration
}

func GetDefaultConfig() Config {
	return Config{DequeueWaitTime: 5 * time.Second}
}
