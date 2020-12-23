package config

import (
	"expvar"
	"time"
)

// Value is the config value interface.
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

// Source is a config source.
type Source interface {
	Resolve(key string) (Value, error)
}

// Watcher is a config watcher.
type Watcher interface {
	Next() (Value, error)
}

// Config is an interface abstraction for configuration.
type Config interface {
	Var(v expvar.Var) error
	Value(key string) Value
	Watch(key ...string) (Watcher, error)
}
