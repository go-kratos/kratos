package queue

// PublishOption is a publish option.
type PublishOption func(*PublishOptions)

// PublishOptions is a publish options.
type PublishOptions struct{}

// NewPublishOptions new a default publish options.
func NewPublishOptions() PublishOptions {
	return PublishOptions{}
}

// SubscribeOption is a subscribe option.
type SubscribeOption func(*SubscribeOptions)

// SubscribeOptions is a subcribe options.
type SubscribeOptions struct {
	AutoAck bool
}

// NewSubscribeOptions new a default subscribe options.
func NewSubscribeOptions() SubscribeOptions {
	return SubscribeOptions{
		AutoAck: true,
	}
}

// DisableAutoAck returns a SubscribeOption which disables auto ack for this
// Subscriber.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}
