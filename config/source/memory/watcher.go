package memory

import "github.com/go-kratos/kratos/v2/config/source"

type watcher struct {
	ch chan *source.KeyValue
}

func newWatcher() source.Watcher {
	return &watcher{ch: make(chan *source.KeyValue)}
}

func (w *watcher) Next() (*source.KeyValue, error) {
	return <-w.ch, nil
}

func (w *watcher) Close() error {
	return nil
}
