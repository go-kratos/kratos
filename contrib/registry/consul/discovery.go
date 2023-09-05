package consul

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

type kratosDiscovery struct {
	// native consul client
	cli *api.Client

	// a main context for monitor goroutine
	ctx context.Context

	resolver    ServiceResolver
	oldBehavior bool

	registry map[string]*serviceSet
	lock     sync.RWMutex
}

type DiscoveryOption func(*kratosDiscovery)

// ServiceResolver is used to resolve service endpoints
type ServiceResolver func(ctx context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance

// WithServiceResolver with endpoint function option.
func WithServiceResolver(fn ServiceResolver) DiscoveryOption {
	return func(o *kratosDiscovery) {
		o.resolver = fn
	}
}

// GetService will query service from consul if not exist in cache
// However, per interface definition, this should be an "in-memory" process
// Here by default, we will only return the cached service instance
// If your service is not watched before, calling GetService without old behavior
// enabled will results empty
func WithOldGetServiceBehavior() DiscoveryOption {
	return func(o *kratosDiscovery) {
		o.oldBehavior = true
	}
}

// @deprecated useless now, timeout is hardcoded, remove this option in your code
func WithTimeout(timeout time.Duration) DiscoveryOption {
	return func(o *kratosDiscovery) {
	}
}

const CONSUL_WAIT_TIME = 55 * time.Second
const CONSUL_API_CONTEXT_TIMEOUT = 60 * time.Second

func NewDiscovery(ctx context.Context, apiClient *api.Client, opts ...DiscoveryOption) registry.Discovery {
	d := &kratosDiscovery{
		ctx:      ctx,
		cli:      apiClient,
		resolver: defaultResolver,
		registry: make(map[string]*serviceSet),
	}

	for _, o := range opts {
		o(d)
	}

	return d
}

func (d *kratosDiscovery) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	d.lock.RLock()
	set := d.registry[name]
	d.lock.RUnlock()

	var ss map[string][]*registry.ServiceInstance

	getRemote := func() map[string][]*registry.ServiceInstance {
		services, _, _, err := d.queryService(ctx, name, map[string]uint64{}, true)
		if err == nil {
			return nil
		}
		return services
	}

	if set == nil || set.services == nil {
		if d.oldBehavior {
			if rss := getRemote(); rss != nil {
				ss = rss
			} else {
				return nil, fmt.Errorf("fail to fetch service %s from consul", name)
			}
		} else {
			return nil, fmt.Errorf("service %s not watched before, no data", name)
		}
	} else {
		ss = set.services
	}

	if ss == nil {
		return nil, fmt.Errorf("service %s no data", name)
	}

	flatSs := make([]*registry.ServiceInstance, 0)

	set.lock.RLock()
	for _, s := range ss {
		flatSs = append(flatSs, s...)
	}
	set.lock.RUnlock()
	return flatSs, nil
}

func (d *kratosDiscovery) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	set, ok := d.registry[name]
	if !ok {
		set = &serviceSet{
			watcher:     make(map[*watcher]struct{}),
			services:    make(map[string][]*registry.ServiceInstance),
			serviceName: name,
		}
		d.registry[name] = set
	}

	// init watcher
	w := &watcher{
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.set = set

	set.lock.Lock()
	set.watcher[w] = struct{}{}
	ss := set.flatServices()
	set.lock.Unlock()

	if len(ss) > 0 {
		// If the service has a value, it needs to be pushed to the watcher,
		// otherwise the initial data may be blocked forever during the watch.
		w.event <- struct{}{}
	}

	if !ok {
		d.monitorUpdate(set)
	}
	return w, nil
}

// return service list.
func (r *kratosDiscovery) ListServices(ctx context.Context) (allServices map[string][]*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	allServices = make(map[string][]*registry.ServiceInstance)
	for name := range r.registry {
		allServices[name], err = r.GetService(ctx, name)
		if err != nil {
			return nil, err
		}
	}
	return
}

// Query consul and return the service instances with index for next call
func (d *kratosDiscovery) queryService(ctx context.Context, service string, indexs map[string]uint64, passingOnly bool) (map[string][]*registry.ServiceInstance, map[string]uint64, bool, error) {
	var instancesByDc = map[string][]*registry.ServiceInstance{}
	var newIndexs = map[string]uint64{}
	var wg sync.WaitGroup
	var anyErr error
	var lock sync.Mutex
	anyUpdated := false

	dcs, err := d.cli.Catalog().Datacenters()
	if err != nil {
		return nil, nil, false, err
	}

	queryCtx, cancelQuery := context.WithCancel(ctx)

	for _, dc := range dcs {
		opts := &api.QueryOptions{
			WaitIndex: indexs[dc],
			WaitTime:  CONSUL_WAIT_TIME,
		}
		opts = opts.WithContext(queryCtx)
		opts.Datacenter = dc

		wg.Add(1)

		// for use in async go routine clousure
		localDc := dc
		go func() {
			defer wg.Done()
			e, m, err := d.cli.Health().Service(service, "", passingOnly, opts)
			useCachedInsts := false
			newIndex := indexs[localDc]
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					anyErr = err
					return
				} else {
					useCachedInsts = true
				}
			} else {
				if m.LastIndex == indexs[localDc] {
					useCachedInsts = true
				} else {
					anyUpdated = true
					newIndex = m.LastIndex
					// others might be blocking, tell them cancel now but with a small delay
					time.AfterFunc(200*time.Millisecond, func() {
						cancelQuery()
					})
				}
			}

			var ins []*registry.ServiceInstance
			if useCachedInsts {
				d.lock.RLock()
				set := d.registry[service]
				d.lock.RUnlock()

				if set != nil {
					set.lock.Lock()
					ss := set.services[localDc]
					set.lock.Unlock()

					if ss != nil {
						ins = ss
					}
				}
			} else {
				ins = d.resolver(ctx, e)
				// for _, in := range ins {
				// 	if in.Metadata == nil {
				// 		in.Metadata = make(map[string]string, 1)
				// 	}
				// 	in.Metadata["dc"] = localDc
				// }
			}

			lock.Lock()
			instancesByDc[localDc] = ins
			newIndexs[localDc] = newIndex
			lock.Unlock()
		}()
	}

	wg.Wait()
	cancelQuery()

	if anyErr != nil {
		return nil, nil, false, anyErr
	}
	return instancesByDc, newIndexs, anyUpdated, nil
}

func defaultResolver(_ context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance {
	services := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		version, vex := entry.Service.Meta["version"]

		if vex {
			// version is a special key, remove it directly and remove whole meta if no others field
			if len(entry.Service.Meta) == 1 {
				entry.Service.Meta = nil
			} else {
				delete(entry.Service.Meta, "version")
			}
		}

		endpoints := make([]string, 0)
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			endpoints = append(endpoints, addr.Address)
		}
		if len(endpoints) == 0 && entry.Service.Address != "" && entry.Service.Port != 0 {
			endpoints = append(endpoints, fmt.Sprintf("http://%s:%d", entry.Service.Address, entry.Service.Port))
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: endpoints,
		})
	}

	return services
}

// monitorUpdate will monitor the service update and notify all watchers
// it will perform initial query first, then start a goroutine to monitor the update
func (d *kratosDiscovery) monitorUpdate(ss *serviceSet) {
	idxs := map[string]uint64{}

	updateBroadcast := func() {
		timeoutCtx, cancel := context.WithTimeout(d.ctx, CONSUL_API_CONTEXT_TIMEOUT)
		tmpService, tmpIdxs, updated, err := d.queryService(timeoutCtx, ss.serviceName, idxs, true)
		fmt.Println("???", ss.serviceName)
		cancel()
		if err != nil {
			return
		}

		if updated {
			ss.broadcast(tmpService)
		}
		idxs = tmpIdxs
	}

	updateBroadcast()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updateBroadcast()
			case <-d.ctx.Done():
				return
			}
		}
	}()
}
