package config

// KeyValue is config key value.
type KeyValue struct {
	Key    string
	Format string
	Value  []byte
}

// Source is config source.
type Source interface {
	Load() ([]*KeyValue, error)
	Watch() (Watcher, error)
}

// Watcher watches a source for changes.
type Watcher interface {
	Next() ([]*KeyValue, error)
	Stop() error
}
