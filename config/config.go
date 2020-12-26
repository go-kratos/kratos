package config

import (
	"expvar"

	"github.com/go-kratos/kratos/v2/config/provider"
)

// Config is a config interface.
type Config interface {
	Var(key string, v expvar.Var) error
	Value(key string) Value
	Watch(key ...string) (Watcher, error)
}

// Option is config option.
type Option func(*options)

type options struct {
	providers []provider.Provider
}

type config struct {
	vars      map[string][]expvar.Var
	resolvers []Resolver
}

func (c *config) setVar(key string, v expvar.Var) error {
	val := c.Value(key)
	switch vv := v.(type) {
	case *expvar.Int:
		intVal, err := val.Int64()
		if err != nil {
			return err
		}
		vv.Set(intVal)
	case *expvar.Float:
		floatVal, err := val.Float64()
		if err != nil {
			return err
		}
		vv.Set(floatVal)
	case *expvar.String:
		stringVal, err := val.String()
		if err != nil {
			return err
		}
		vv.Set(stringVal)
	}
	return nil
}

func (c *config) Var(key string, v expvar.Var) error {
	if err := c.setVar(key, v); err != nil {
		return err
	}
	c.vars[key] = append(c.vars[key], v)
	return nil
}

func (c *config) Value(key string) Value {
	for _, r := range c.resolvers {
		v, ok := r.Resolve(key)
		if ok {
			return v
		}
	}
	return &errValue{err: ErrNotFound}
}

func (c *config) Watch(key ...string) (Watcher, error) {
	// TODO
	return nil, nil
}
