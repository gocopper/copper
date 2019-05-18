package cpubsub

type Handler func(payload []byte)

type PubSub interface {
	Subscribe(topic string, handler *Handler) error
	Publish(topic string, payload []byte) error
}
