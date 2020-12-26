package config

// Watcher is config watcher.
type Watcher interface {
	Next() (Value, error)
}
