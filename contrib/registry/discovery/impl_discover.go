package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/go-kratos/kratos/v2/registry"
)

func filterInstancesByZone(ins *disInstancesInfo, zone string) []*registry.ServiceInstance {
	zoneInstance, ok := ins.Instances[zone]
	if !ok || len(zoneInstance) == 0 {
		return nil
	}

	out := make([]*registry.ServiceInstance, 0, len(zoneInstance))
	for _, v := range zoneInstance {
		if v == nil {
			continue
		}
		out = append(out, toServiceInstance(v))
	}

	return out
}

func (d *Discovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	r := d.resolveBuild(serviceName)
	ins, ok := r.fetch(ctx)
	if !ok {
		return nil, errors.New("Discovery.GetService fetch failed")
	}

	out := filterInstancesByZone(ins, d.config.Zone)
	if len(out) == 0 {
		return nil, fmt.Errorf("Discovery.GetService(%s) not found", serviceName)
	}

	return out, nil
}

func (d *Discovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &watcher{
		Resolve:     d.resolveBuild(serviceName),
		serviceName: serviceName,
		cancelCtx:   ctx,
	}, nil
}

type watcher struct {
	*Resolve

	cancelCtx   context.Context
	serviceName string
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	event := w.Resolve.Watch()

	select {
	case <-event:
	// change event come
	case <-w.cancelCtx.Done():
		return nil, fmt.Errorf("watch context cancelled: %v", w.cancelCtx.Err())
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	ins, ok := w.Resolve.fetch(ctx)
	if !ok {
		return nil, errors.New("Discovery.GetService fetch failed")
	}

	out := filterInstancesByZone(ins, w.Resolve.d.config.Zone)
	if len(out) == 0 {
		return nil, fmt.Errorf("Discovery.GetService(%s) not found", w.serviceName)
	}

	return out, nil
}

func (w *watcher) Stop() error {
	return w.Resolve.Close()
}
