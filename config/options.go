package config

import (
	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/source"
)

// Option is config option.
type Option func(*options)

type options struct {
	parsers []parser.Parser
	sources []source.Source
}

// WithSource .
func WithSource(s ...source.Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

// WithParser .
func WithParser(p ...parser.Parser) Option {
	return func(o *options) {
		o.parsers = p
	}
}
