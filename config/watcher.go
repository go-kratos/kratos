package config

// Watcher is config watcher.
type Watcher interface {
	Next() (Value, error)
	Close() error
}

type watcher struct {
	ch chan Value
}

func newWatcher() Watcher {
	return &watcher{}
}

func (w *watcher) update(v Value) {
	w.ch <- v
}

func (w *watcher) Next() (Value, error) {
	return <-w.ch, nil
}

func (w *watcher) Close() error { return nil }
