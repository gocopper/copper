package cpubsub

import "go.uber.org/fx"

var RedisFx = fx.Provide(
	NewRedisPubSub,
)
