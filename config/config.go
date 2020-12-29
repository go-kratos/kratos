package config

import (
	"errors"
	"expvar"
	"sync"

	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/parser/json"
	"github.com/go-kratos/kratos/v2/config/parser/toml"
	"github.com/go-kratos/kratos/v2/config/parser/yaml"
)

var (
	// ErrNotFound is value not found.
	ErrNotFound = errors.New("error key not found")

	_ Config = (*config)(nil)
)

// Config is a config interface.
type Config interface {
	Load() error
	Var(key string, v expvar.Var) error
	Value(key string) Value
	Watch(key ...string) <-chan Value
}

type config struct {
	vars      sync.Map
	cached    sync.Map
	resolvers []Resolver
	opts      options
}

// New new a config with options.
func New(opts ...Option) Config {
	options := options{
		parsers: []parser.Parser{
			json.NewParser(),
			yaml.NewParser(),
			toml.NewParser(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return &config{opts: options}
}

func (c *config) Load() error {
	parsers := make(map[string]parser.Parser)
	for _, parser := range c.opts.parsers {
		parsers[parser.Format()] = parser
	}
	for _, p := range c.opts.providers {
		r, err := newResolver(p, parsers)
		if err != nil {
			return err
		}
		c.resolvers = append(c.resolvers, r)
	}
	return nil
}

func (c *config) Var(key string, v expvar.Var) error {
	if err := setVar(v, c.Value(key)); err != nil {
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
