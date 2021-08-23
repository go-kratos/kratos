package registry

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

// Option is consul registry option.
type Option func(*Registry)

// WithHealthCheck with registry health check option.
func WithHealthCheck(enable bool) Option {
	return func(o *Registry) {
		o.enableHealthCheck = enable
	}
}

// Config is consul registry config
type Config struct {
	*api.Config
}

// Registry is consul registry
type Registry struct {
	cfg               *Config
	cli               *Client
	enableHealthCheck bool

	registry map[string]*serviceSet
	lock     sync.RWMutex
}

// New creates consul registry
func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{
		cli:               NewClient(apiClient),
		registry:          make(map[string]*serviceSet),
		enableHealthCheck: true,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Register register service
func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

// Deregister deregister service
func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

// GetService return service by name
func (r *Registry) GetService(ctx context.Context, name string) (services []*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	set := r.registry[name]
	if set == nil {
		return nil, fmt.Errorf("service %s not resolved in registry", name)
	}
	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if ss == nil {
		return nil, fmt.Errorf("service %s not found in registry", name)
	}
	for _, s := range ss {
		services = append(services, s)
	}
	return
}

// ListServices return service list.
func (r *Registry) ListServices() (allServices map[string][]*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	allServices = make(map[string][]*registry.ServiceInstance)
	for name, set := range r.registry {
		var services []*registry.ServiceInstance
		ss, _ := set.services.Load().([]*registry.ServiceInstance)
		if ss == nil {
			continue
		}
		for _, s := range ss {
			services = append(services, s)
		}
		allServices[name] = services
	}
	return
}

// Watch resolve service by name
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	set, ok := r.registry[name]
	if !ok {
		set = &serviceSet{
			watcher:     make(map[*watcher]struct{}, 0),
			services:    &atomic.Value{},
			serviceName: name,
		}
		r.registry[name] = set
	}

	// 初始化watcher
	w := &watcher{
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.set = set
	set.lock.Lock()
	set.watcher[w] = struct{}{}
	set.lock.Unlock()
	ss, _ := set.services.Load().([]*registry.ServiceInstance)
	if len(ss) > 0 {
		// 如果services有值需要推送给watcher，否则watch的时候可能会永远阻塞拿不到初始的数据
		w.event <- struct{}{}
	}

	// 放在最后是为了防止漏推送
	if !ok {
		go r.resolve(set)
	}
	return w, nil
}

func (r *Registry) resolve(ss *serviceSet) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	services, idx, err := r.cli.Service(ctx, ss.serviceName, 0, true)
	cancel()
	if err == nil && len(services) > 0 {
		ss.broadcast(services)
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		tmpService, tmpIdx, err := r.cli.Service(ctx, ss.serviceName, idx, true)
		cancel()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if len(tmpService) != 0 && tmpIdx != idx {
			services = tmpService
			ss.broadcast(services)
		}
		idx = tmpIdx
	}
}
