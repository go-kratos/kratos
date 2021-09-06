package zap

import (
	"go.uber.org/zap"
)

type options struct {
	zapConfig  zap.Config
	zapOptions []zap.Option
}

type Option func(*options)

func WithConfig(config zap.Config) Option {
	return func(o *options) {
		o.zapConfig = config
	}
}

func WithOptions(opts ...zap.Option) Option {
	return func(o *options) {
		o.zapOptions = opts
	}
}
