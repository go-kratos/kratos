package consul

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type watcher struct {
	r  *Registry
	ch chan string
	wp *watch.Plan

	// for cancel
	ctx    context.Context
	cancel context.CancelFunc
}

func newWatcher(ctx context.Context, serviceName string, r *Registry) (*watcher, error) {
	ctx, cancel := context.WithCancel(ctx)
	w := &watcher{
		r:  r,
		ch: make(chan string),

		ctx:    ctx,
		cancel: cancel,
	}

	wp, err := watch.Parse(map[string]interface{}{
		"type":    "service",
		"service": serviceName,
	})
	if err != nil {
		return nil, err
	}

	wp.Handler = w.handle
	w.wp = wp

	// wp.Run is a blocking call and will prevent newWatcher from returning
	go func() {
		err := wp.RunWithClientAndHclog(r.cli.cli, nil)
		if err != nil {
			panic(err)
		}
	}()

	return w, nil
}

func (w *watcher) handle(idx uint64, data interface{}) {
	if data == nil {
		return
	}

	m := make(map[string]struct{})
	switch d := data.(type) {
	case []*api.ServiceEntry:
		for _, i := range d {
			if i.Checks.AggregatedStatus() == api.HealthPassing {
				m[i.Service.Service] = struct{}{}
			}
		}
	}

	for name := range m {
		w.ch <- name
	}
}

func (w *watcher) Next() (services []*registry.ServiceInstance, err error) {
	select {
	case name := <-w.ch:
		return w.r.GetService(context.Background(), name)
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	w.wp.Stop()
	w.cancel()
	return nil
}
