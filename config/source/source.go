package source

import "time"

// KeyValue is config key value.
// format: json/yaml/text
type KeyValue struct {
	Key       string
	Value     []byte
	Format    string
	Timestamp time.Time
}

// Source is config source.
type Source interface {
	Load() ([]*KeyValue, error)
	Watch() (Watcher, error)
}

// Watcher watches a provider for changes
type Watcher interface {
	Next() (*KeyValue, error)
	Close() error
}
