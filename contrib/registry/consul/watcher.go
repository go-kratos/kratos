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

	ss, ok := w.set.services.Load().([]*registry.ServiceInstance)

	if ok {
		services = append(services, ss...)
	}
	return
}

func (w *watcher) Stop() error {
	w.cancel()
	w.set.lock.Lock()
	defer w.set.lock.Unlock()
	delete(w.set.watcher, w)
	return nil
}
