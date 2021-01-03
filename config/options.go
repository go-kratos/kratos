package config

import (
	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/parser/json"
	"github.com/go-kratos/kratos/v2/config/parser/text"
	"github.com/go-kratos/kratos/v2/config/parser/toml"
	"github.com/go-kratos/kratos/v2/config/parser/yaml"
	"github.com/go-kratos/kratos/v2/config/source"
)

// Option is config option.
type Option func(*options)

type options struct {
	parsers []parser.Parser
	sources []source.Source
}

func defaultOptions() options {
	return options{
		parsers: []parser.Parser{
			text.NewParser(),
			json.NewParser(),
			yaml.NewParser(),
			toml.NewParser(),
		},
	}
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
