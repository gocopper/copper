package cpubsub

// LocalPubSub is an implementation of PubSub that only allows publishing
// of payloads to 'local' subscribers i.e. subscribers on the same instance
// as the publisher.
type LocalPubSub struct {
	subscriptions map[string][]*Handler
}

func NewLocalPubSub() *LocalPubSub {
	return &LocalPubSub{
		subscriptions: make(map[string][]*Handler),
	}
}

func (l *LocalPubSub) Subscribe(topic string, handler *Handler) error {
	handlers, ok := l.subscriptions[topic]
	if !ok {
		l.subscriptions[topic] = []*Handler{handler}
		return nil
	}

	l.subscriptions[topic] = append(handlers, handler)

	return nil
}

func (l *LocalPubSub) Publish(topic string, payload []byte) error {
	handlers, ok := l.subscriptions[topic]
	if !ok {
		return nil
	}

	for _, h := range handlers {
		go (*h)(payload)
	}

	return nil
}
