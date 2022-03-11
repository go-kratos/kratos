package env

import (
	"context"

	"github.com/SeeMusic/kratos/v2/config"
)

type watcher struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var _ config.Watcher = (*watcher)(nil)

func NewWatcher() (config.Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{ctx: ctx, cancel: cancel}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	<-w.ctx.Done()
	return nil, w.ctx.Err()
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}
