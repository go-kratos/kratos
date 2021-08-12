package env

import (
	"context"

	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	exit chan struct{}

	ctx    context.Context
	cancel context.CancelFunc
}

func NewWatcher() (config.Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{exit: make(chan struct{}), ctx: ctx, cancel: cancel}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.exit:
		return nil, nil
	}
}

func (w *watcher) Stop() error {
	close(w.exit)
	w.cancel()
	return nil
}
