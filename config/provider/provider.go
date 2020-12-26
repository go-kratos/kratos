package provider

import "time"

// KeyValue is config key value.
// format: json/yaml/text
type KeyValue struct {
	Key       string
	Value     []byte
	Format    string
	Timestamp time.Time
}

// Watcher is config watcher.
type Watcher interface {
	Next() (KeyValue, error)
}

// Provider is config provider.
type Provider interface {
	Load() ([]KeyValue, error)
	Watch() (Watcher, error)
}
