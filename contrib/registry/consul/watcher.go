package consul

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
)

type watcher struct {
	event chan struct{}
	set   *serviceSet

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *watcher) Next() (services []*registry.ServiceInstance, err error) {
	select {
	case <-w.ctx.Done():
		err = w.ctx.Err()
		return
	case <-w.event:
	}

	if err = w.ctx.Err(); err != nil {
		return nil, err
	}

	return w.set.getInstances(), nil
}

func (w *watcher) Stop() error {
	w.cancel()
	w.set.delWatcher(w)
	return nil
}
