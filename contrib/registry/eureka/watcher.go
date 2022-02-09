package eureka

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
)

var _ registry.Watcher = &watcher{}

type watcher struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cli        *API
	watchChan  chan struct{}
	serverName string
}

func newWatch(ctx context.Context, cli *API, serverName string) (*watcher, error) {
	w := &watcher{
		ctx:        ctx,
		cli:        cli,
		serverName: serverName,
		watchChan:  make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	e := w.cli.Subscribe(
		serverName,
		func() {
			w.watchChan <- struct{}{}
		},
	)
	return w, e
}

func (w watcher) Next() (services []*registry.ServiceInstance, err error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
		instances := w.cli.GetService(w.ctx, w.serverName)
		services = make([]*registry.ServiceInstance, 0, len(instances))
		for _, instance := range instances {
			services = append(services, &registry.ServiceInstance{
				ID:        instance.Metadata["ID"],
				Name:      instance.Metadata["Name"],
				Version:   instance.Metadata["Version"],
				Endpoints: []string{instance.Metadata["Endpoints"]},
				Metadata:  instance.Metadata,
			})
		}
		return
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	w.cli.Unsubscribe(w.serverName)
	return nil
}
