package config

import (
	"expvar"
	"sync"

	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/provider"
)

// Config is a config interface.
type Config interface {
	Var(key string, v expvar.Var) error
	Value(key string) Value
	Watch(key ...string) <-chan Value
}

// Option is config option.
type Option func(*options)

type options struct {
	providers []provider.Provider
	parsers   map[string]parser.Parser
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
		if o.parsers == nil {
			o.parsers = make(map[string]parser.Parser)
		}
		for _, parser := range p {
			o.parsers[parser.Format()] = parser
		}
	}
}

type config struct {
	vars      sync.Map
	cached    sync.Map
	resolvers []Resolver
	opts      options
}

// New new a config with options.
func New(opts ...Option) Config {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	c := &config{opts: options}
	for _, p := range options.providers {
		c.resolvers = append(c.resolvers, newResolver(p, options.parsers))
	}
	return c
}

func (c *config) Var(key string, v expvar.Var) error {
	if err := setVar(key, v, c.Value(key)); err != nil {
		return err
	}
	c.vars.Store(key, v)
	return nil
}

func (c *config) Value(key string) Value {
	if v, ok := c.cached.Load(key); ok {
		return v.(Value)
	}
	for _, r := range c.resolvers {
		v, ok := r.Resolve(key)
		if ok {
			c.cached.Store(key, v)
			return v
		}
	}
	return &errValue{err: ErrNotFound}
}

func (c *config) Watch(key ...string) <-chan Value {
	// TODO
	return nil
}
