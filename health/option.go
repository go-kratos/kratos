package health

import "time"

type Option func(o *Health)

func WithTimeout(timeout time.Duration) Option {
	return func(o *Health) {
		o.timeout = timeout
	}
}

func WithIntervalTime(intervalTime time.Duration) Option {
	return func(o *Health) {
		o.intervalTime = intervalTime
	}
}
