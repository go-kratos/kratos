package config

import (
	"github.com/go-kratos/kratos/v2/log"
)

// Decoder is config decoder.
type Decoder func(*KeyValue, map[string]interface{}) error

// Option is config option.
type Option func(*options)

type options struct {
	sources []Source
	decoder Decoder
	logger  log.Logger
}

// WithSource with config source.
func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

// WithDecoder with config decoder.
func WithDecoder(d Decoder) Option {
	return func(o *options) {
		o.decoder = d
	}
}

// WithLogger with config loogger.
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}
