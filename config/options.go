package config

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Decoder is config decoder.
type Decoder func(*KeyValue, map[string]interface{}) error

// Option is config option.
type Option func(*Options)

type Options struct {
	sources []Source
	decoder Decoder
	logger  log.Logger

	// for alternative data
	Context context.Context
}

// WithSource with config source.
func WithSource(s ...Source) Option {
	return func(o *Options) {
		o.sources = s
	}
}

// WithDecoder with config decoder.
func WithDecoder(d Decoder) Option {
	return func(o *Options) {
		o.decoder = d
	}
}

// WithLogger with config loogger.
func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}
