package config

import "github.com/go-kratos/kratos/v2/config/source"

// Watcher is config watcher.
type Watcher interface {
	Next() (Value, error)
	Close() error
}

// Observer is config watch observer.
type Observer func(*source.KeyValue)

type watcher struct {
}

func newWatcher() Watcher {
	return &watcher{}
}

func (w *watcher) Next() (Value, error) {
	return nil, nil
}

func (w *watcher) Close() error { return nil }
