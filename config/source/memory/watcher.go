package memory

import "github.com/go-kratos/kratos/v2/config/source"

type watcher struct {
	ch chan *source.KeyValue
}

func newWatcher(ch chan *source.KeyValue) source.Watcher {
	return &watcher{ch: ch}
}

func (w *watcher) Next() (*source.KeyValue, error) {
	return <-w.ch, nil
}

func (w *watcher) Close() error {
	return nil
}
