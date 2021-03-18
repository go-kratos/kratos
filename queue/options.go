package queue

// PublishOption is a publish option.
type PublishOption func(*PublishOptions)

// PublishOptions is a publish options.
type PublishOptions struct {
	Header map[string]string
}

// WithHeader with a customized user header delivering to the message.
func WithHeader(h Header) PublishOption {
	return func(o *PublishOptions) {
		o.Header = h
	}
}

// SubscribeOption is a subscribe option.
type SubscribeOption func(*SubscribeOptions)

// SubscribeOptions is a subcribe options.
type SubscribeOptions struct {
	AutoAck bool
}

// DisableAutoAck returns a SubscribeOption which disables auto ack for this
// Subscriber.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}
