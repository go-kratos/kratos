package metrics

import "github.com/DataDog/datadog-go/statsd"

// Option is doatadog option.
type Option func(*options)

type options struct {
	sampleRate float64
	labels     []string
	client     *statsd.Client
}

// WithSampleRate with sample rate option.
func WithSampleRate(rate float64) Option {
	return func(o *options) {
		o.sampleRate = rate
	}
}

// WithLabels with labels option.
func WithLabels(lvs ...string) Option {
	return func(o *options) {
		o.labels = lvs
	}
}

// WithClient with client option.
func WithClient(c *statsd.Client) Option {
	return func(o *options) {
		o.client = c
	}
}
