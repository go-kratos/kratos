package queue

// PublishOption is a publish option.
type PublishOption func(*PublishOptions)

// PublishOptions is a publish options.
type PublishOptions struct{}

// SubscribeOption is a subscribe option.
type SubscribeOption func(*SubscribeOptions)

// SubscribeOptions is a subcribe options.
type SubscribeOptions struct{}
