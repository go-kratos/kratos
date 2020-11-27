package config

import (
	"expvar"
	"time"
)

// MD is the value metadata.
type MD interface {
	Key() string
}

// Value is the config value interface.
type Value interface {
	Metadata() MD

	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration

	Scan(val interface{}) error
}

// Watcher is the config watcher.
type Watcher interface {
	Next() (Value, error)
}

// Config is an interface abstraction for configuration.
type Config interface {
	Var(v expvar.Var) error
	Value(key string) Value
	Watch(key ...string) (Watcher, error)
}
