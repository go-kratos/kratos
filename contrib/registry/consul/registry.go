package consul

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

var (
	_ registry.Registrar = (*Registry)(nil)
	_ registry.Discovery = (*Registry)(nil)
)

// Option is consul registry option.
type Option func(*Registry)

// WithHealthCheck with registry health check option.
func WithHealthCheck(enable bool) Option {
	return func(o *Registry) {
		o.enableHealthCheck = enable
	}
}

// WithTimeout with get services timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(o *Registry) {
		o.timeout = timeout
	}
}

// WithClusters specify the cluster to be used, if not set, obtain all currently associated clusters.
func WithClusters(cs ...string) Option {
	return func(o *Registry) {
		o.cli.clusters = cs
	}
}

func WithMultiClusterMode(mode ClusterMode) Option {
	return func(r *Registry) {
		r.cli.multiClusterMode = mode
	}
}

// WithHeartbeat enable or disable heartbeat
func WithHeartbeat(enable bool) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.heartbeat = enable
		}
	}
}

// WithServiceResolver with endpoint function option.
func WithServiceResolver(fn ServiceResolver) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.resolver = fn
		}
	}
}

// WithHealthCheckInterval with healthcheck interval in seconds.
func WithHealthCheckInterval(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.healthcheckInterval = interval
		}
	}
}

// WithDeregisterCriticalServiceAfter with deregister-critical-service-after in seconds.
func WithDeregisterCriticalServiceAfter(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.deregisterCriticalServiceAfter = interval
		}
	}
}

// WithServiceCheck with service checks
func WithServiceCheck(checks ...*api.AgentServiceCheck) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.serviceChecks = checks
		}
	}
}

// WithAllowReRegistration allow re registration when the service is automatically removed by the registry
func WithAllowReRegistration(allow bool) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.allowReRegistration = allow
		}
	}
}

// Config is consul registry config
type Config struct {
	*api.Config
}

// Registry is consul registry
type Registry struct {
	cli               *Client
	enableHealthCheck bool
	registry          map[string]*serviceSet
	lock              sync.RWMutex
	timeout           time.Duration
}

// New creates consul registry
func New(apiClient *api.Client, opts ...Option) *Registry {
	r := &Registry{
		registry:          make(map[string]*serviceSet),
		enableHealthCheck: true,
		timeout:           10 * time.Second,
		cli: &Client{
			consul:                         apiClient,
			resolver:                       defaultResolver,
			healthcheckInterval:            10,
			heartbeat:                      true,
			deregisterCriticalServiceAfter: 600,
			multiClusterMode:               Single,
			deregisteredService:            make(map[string]struct{}),
		},
	}
	for _, o := range opts {
		o(r)
	}

	var err error
	if r.cli.multiClusterMode == WanFederation || r.cli.multiClusterMode == Peering {
		if len(r.cli.clusters) == 0 {
			switch r.cli.multiClusterMode {
			case WanFederation:
				r.cli.clusters, err = r.cli.consul.Catalog().Datacenters()
				if err != nil {
					log.Warnf("[Consul] get datacenters failed，will use the current cluster! err=%v", err)
				}
			case Peering:
				peerings, _, err := r.cli.consul.Peerings().List(context.Background(), nil)
				if err != nil {
					log.Warnf("[Consul] get peerings failed，will use the current cluster! err=%v", err)
					break
				}
				for _, peering := range peerings {
					r.cli.clusters = append(r.cli.clusters, peering.Name)
				}
			}
		}
	}

	// append the current cluster
	r.cli.clusters = append(r.cli.clusters, "")

	r.cli.ctx, r.cli.cancel = context.WithCancel(context.Background())
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
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	var services []*registry.ServiceInstance
	for _, cluster := range r.cli.clusters {
		opts := &api.QueryOptions{}
		switch r.cli.multiClusterMode {
		case WanFederation:
			opts.Datacenter = cluster
		case Peering:
			opts.Peer = cluster
		}

		tmp, _, err := r.cli.service(ctx, name, true, opts)
		if err != nil {
			return nil, err
		}
		services = append(services, tmp...)
	}

	if len(services) == 0 {
		return nil, errors.New("service instances not found")
	}

	return services, nil
}

// Deprecated: This method should not be exported
// ListServices return service list.
func (r *Registry) ListServices() (allServices map[string][]*registry.ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	allServices = make(map[string][]*registry.ServiceInstance)
	for name, set := range r.registry {
		ss := set.getInstances()
		if ss == nil {
			continue
		}
		allServices[name] = ss
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
			watcher:     make(map[*watcher]struct{}),
			services:    make(map[string][]*registry.ServiceInstance),
			serviceName: name,
		}
		r.registry[name] = set
	}

	// init watcher
	w := &watcher{
		event: make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.set = set
	set.addWatcher(w)
	ss := set.getInstances()
	if len(ss) > 0 {
		// If the service has a value, it needs to be pushed to the watcher,
		// otherwise the initial data may be blocked forever during the watch.
		w.event <- struct{}{}
	}

	if !ok {
		err := r.resolve(ctx, set)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (r *Registry) resolve(ctx context.Context, ss *serviceSet) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	for _, cluster := range r.cli.clusters {
		res, _, err := r.cli.service(timeoutCtx, ss.serviceName, true, processClusterOption(nil, r.cli.multiClusterMode, cluster))
		if err != nil {
			return err
		}
		if len(res) > 0 {
			ss.broadcast(cluster, res)
		}
	}

	for _, cluster := range r.cli.clusters {
		go func(cluster string) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			opts := processClusterOption(nil, r.cli.multiClusterMode, cluster)
			opts.WaitIndex = 0
			opts.WaitTime = time.Second * 55

			var err error
			var tmpService []*registry.ServiceInstance

			for range ticker.C {
				timeoutCtx, cancel := context.WithTimeout(context.Background(), r.timeout)
				tmpService, opts.WaitIndex, err = r.cli.service(timeoutCtx, ss.serviceName, true, opts)
				cancel()
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				if len(tmpService) > 0 {
					ss.broadcast(cluster, tmpService)
				}
			}
		}(cluster)
	}
	return nil
}

func processClusterOption(opts *api.QueryOptions, mode ClusterMode, cluster string) *api.QueryOptions {
	if opts == nil {
		opts = &api.QueryOptions{}
	}
	switch mode {
	case WanFederation:
		opts.Datacenter = cluster
	case Peering:
		opts.Peer = cluster
	}
	return opts
}
