package config

import (
	"errors"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/parser/json"
	"github.com/go-kratos/kratos/v2/config/parser/toml"
	"github.com/go-kratos/kratos/v2/config/parser/yaml"
	"github.com/go-kratos/kratos/v2/config/source"
)

var (
	// ErrNotFound is key not found.
	ErrNotFound = errors.New("key not found")
	// ErrTypeAssert is type assert error.
	ErrTypeAssert = errors.New("type assert error")

	_ Config = (*config)(nil)
)

// Config is a config interface.
type Config interface {
	Load() error
	Value(key string) Value
	Watch(key ...string) (Watcher, error)
}

type config struct {
	cached    sync.Map
	resolvers []*resolver
	watchers  []source.Watcher
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

func (c *config) watch(r *resolver, w source.Watcher) {
	for {
		kv, err := w.Next()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		r.reload(kv)
		c.cached.Range(func(key, value interface{}) bool {
			for _, r := range c.resolvers {
				if v := r.Resolve(key.(string)); v != value {
					c.cached.Store(key, v)
				}
			}
			return true
		})
	}
}

func (c *config) Load() error {
	parsers := make(map[string]parser.Parser)
	for _, parser := range c.opts.parsers {
		parsers[parser.Format()] = parser
	}
	for _, source := range c.opts.sources {
		w, err := source.Watch()
		if err != nil {
			return err
		}
		r, err := newResolver(source, parsers)
		if err != nil {
			return err
		}
		c.resolvers = append(c.resolvers, r)
		go c.watch(r, w)
	}
	return nil
}

func (c *config) Value(key string) Value {
	if v, ok := c.cached.Load(key); ok {
		return v.(Value)
	}
	for _, r := range c.resolvers {
		if v := r.Resolve(key); v != nil {
			c.cached.Store(key, v)
			return v
		}
	}
	return &errValue{err: ErrNotFound}
}

func (c *config) Watch(key ...string) (Watcher, error) {
	return newWatcher(), nil
}
