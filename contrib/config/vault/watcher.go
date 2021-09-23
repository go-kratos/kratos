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

func (s *watcher) Next() ([]*config.KeyValue, error) {
	<-s.closeChan
	return nil, nil
}

func (s *watcher) Stop() error {
	close(s.closeChan)
	return nil
}
