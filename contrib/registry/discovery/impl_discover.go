package discovery

import (
	"context"
	"fmt"

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

func (d *discovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	r := d.resolveBuild(serviceName)
	ins, ok := r.Fetch(ctx)
	if !ok {
		return nil, errors.New("discovery.GetService fetch failed")
	}

	out := filterInstancesByZone(ins, d.config.Zone)
	if len(out) == 0 {
		return nil, fmt.Errorf("discovery.GetService(%s) not found", serviceName)
	}

	return out, nil
}

func (d *discovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &watcher{
		Resolve:     d.resolveBuild(serviceName),
		serviceName: serviceName,
	}, nil
}

type watcher struct {
	*Resolve

	serviceName string
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	event := w.Resolve.Watch()
	// change event come
	_, ok := <-event

	//ctx, cancel := context.WithTimeout()
	//defer cancel()

	ins, ok := w.Resolve.Fetch(context.TODO())
	if !ok {
		return nil, errors.New("discovery.GetService fetch failed")
	}

	out := filterInstancesByZone(ins, w.Resolve.d.config.Zone)
	if len(out) == 0 {
		return nil, fmt.Errorf("discovery.GetService(%s) not found", w.serviceName)
	}

	return out, nil
}

func (w *watcher) Stop() error {
	return w.Resolve.Close()
}
