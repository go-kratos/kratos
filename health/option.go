package health

import "time"

type Option func(*options)

type options struct {
	watchTime time.Duration
}

func WithWatchTime(t time.Duration) Option {
	return func(o *options) {
		o.watchTime = t
	}
}
