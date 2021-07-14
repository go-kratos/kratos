package env

import (
	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	exit chan struct{}
}

func NewWatcher() (config.Watcher, error) {
	return &watcher{exit: make(chan struct{})}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	<-w.exit
	return nil, nil
}

func (w *watcher) Stop() error {
	close(w.exit)
	return nil
}
