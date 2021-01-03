package config

import (
	"errors"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/config/source"
)

var (
	// ErrNotFound is key not found.
	ErrNotFound = errors.New("key not found")
	// ErrTypeAssert is type assert error.
	ErrTypeAssert = errors.New("type assert error")

	_ Config = (*config)(nil)
)

// Observer is config observer.
type Observer func(string, Value)

// Config is a config interface.
type Config interface {
	Load() error
	Value(key string) Value
	Watch(key string, o Observer) error
	Close() error
}

type config struct {
	cached    sync.Map
	observers sync.Map
	watchers  []source.Watcher
	resolvers []*resolver
	opts      options
}

// New new a config with options.
func New(opts ...Option) Config {
	options := defaultOptions()
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
			k := key.(string)
			v := value.(Value)
			for _, r := range c.resolvers {
				if n := r.Resolve(k); n != nil && n.Load() != v.Load() {
					v.Store(n.Load())
					if o, ok := c.observers.Load(k); ok {
						o.(Observer)(k, v)
					}
				}
			}
			return true
		})
	}
}

func (c *config) Load() error {
	for _, source := range c.opts.sources {
		w, err := source.Watch()
		if err != nil {
			return err
		}
		r, err := newResolver(source, c.opts)
		if err != nil {
			return err
		}
		c.watchers = append(c.watchers, w)
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

func (c *config) Watch(key string, o Observer) error {
	if v := c.Value(key); v.Load() == nil {
		return ErrNotFound
	}
	c.observers.Store(key, o)
	return nil
}

func (c *config) Close() error {
	for _, w := range c.watchers {
		if err := w.Close(); err != nil {
			return err
		}
	}
	return nil
}
