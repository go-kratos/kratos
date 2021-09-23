package vault

import (
	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	source    *source
	closeChan chan struct{}
}

func newWatcher(s *source) (*watcher, error) {
	w := &watcher{
		source:    s,
		closeChan: make(chan struct{}),
	}

	return w, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	<-w.closeChan
	return nil, nil
}

func (w *watcher) Stop() error {
	close(w.closeChan)
	return nil
}
