package zap

import (
	"github.com/go-kratos/kratos/v2/log"
)

type options struct {
	output string
	level  log.Level
	skip   int
	format string
}

type Option func(*options)

func WithOutput(output string) Option {
	return func(o *options) {
		o.output = output
	}
}

func WithLevel(level log.Level) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithSkip(skip int) Option {
	return func(o *options) {
		o.skip = skip
	}
}

func WithFormat(format string) Option {
	return func(o *options) {
		o.format = format
	}
}
