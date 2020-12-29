package config

import (
	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/provider"
)

// Option is config option.
type Option func(*options)

type options struct {
	parsers   []parser.Parser
	providers []provider.Provider
}

// WithProvider .
func WithProvider(p ...provider.Provider) Option {
	return func(o *options) {
		o.providers = p
	}
}

// WithParser .
func WithParser(p ...parser.Parser) Option {
	return func(o *options) {
		o.parsers = p
	}
}
