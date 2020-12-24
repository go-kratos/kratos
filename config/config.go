package config

import (
	"expvar"
	"time"
)

// Value is config value interface.
type Value interface {
	Bool() (bool, error)
	Int() (int, error)
	Int32() (int32, error)
	Int64() (int64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	String() (string, error)
	Duration() (time.Duration, error)
	Scan(interface{}) error
}

// Resolver is config resolver.
type Resolver interface {
	Value(key string) (Value, bool)
}

// Watcher is config watcher.
type Watcher interface {
	Next() (Value, error)
}

// Config is a config interface.
type Config interface {
	Var(v expvar.Var) error
	Value(key string) Value
	Watch(key ...string) (Watcher, error)
}
