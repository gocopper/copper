package cpubsub

import "go.uber.org/fx"

var LocalFx = fx.Provide(
	NewLocalPubSub,
)

var RedisFx = fx.Provide(
	NewRedisPubSub,
)
