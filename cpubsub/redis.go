package cpubsub

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type subscription struct {
	redis    *redis.PubSub
	handlers []*Handler
}

type RedisPubSub struct {
	redis  *redis.Client
	logger clogger.Logger

	subscriptions map[string]*subscription
}

func NewRedisPubSub(config RedisConfig, lc fx.Lifecycle, logger clogger.Logger) PubSub {
	pubsub := &RedisPubSub{
		redis: redis.NewClient(&redis.Options{Addr: config.Addr}),
		logger: logger.WithTags(map[string]interface{}{
			"addr": config.Addr,
		}),
		subscriptions: make(map[string]*subscription),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pubsub.logger.Info("[cpubsub] Connecting to Redis..")

			_, err := pubsub.redis.Ping().Result()
			if err != nil {
				return cerror.New(err, "failed to connect to redis", map[string]interface{}{
					"addr": config.Addr,
				})
			}
			return nil
		},

		OnStop: func(ctx context.Context) error {
			for _, sub := range pubsub.subscriptions {
				pubsub.logger.Info("[cpubsub] Closing Redis pub sub connection..")
				if err := sub.redis.Close(); err != nil {
					return err
				}
			}

			pubsub.logger.Info("[cpubsub] Closing Redis connection..")

			return pubsub.redis.Close()
		},
	})

	return pubsub
}

func (r *RedisPubSub) Subscribe(topic string, handler *Handler) error {
	sub, ok := r.subscriptions[topic]
	if ok {
		sub.handlers = append(sub.handlers, handler)
		return nil
	}

	sub = &subscription{
		redis:    r.redis.Subscribe(topic),
		handlers: []*Handler{handler},
	}

	go func() {
		for m := range sub.redis.Channel() {
			for _, h := range sub.handlers {
				(*h)([]byte(m.Payload))
			}
		}
	}()

	r.subscriptions[topic] = sub

	return nil
}

func (r *RedisPubSub) Publish(topic string, payload []byte) error {
	return r.redis.Publish(topic, payload).Err()
}
