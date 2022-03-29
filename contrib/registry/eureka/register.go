package eureka

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

type Option func(o *Registry)

// WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *Registry) { o.ctx = ctx }
}

func WithHeartbeat(interval time.Duration) Option {
	return func(o *Registry) { o.heartbeatInterval = interval }
}

func WithRefresh(interval time.Duration) Option {
	return func(o *Registry) { o.refreshInterval = interval }
}

func WithEurekaPath(path string) Option {
	return func(o *Registry) { o.eurekaPath = path }
}

type Registry struct {
	ctx               context.Context
	api               *API
	heartbeatInterval time.Duration
	refreshInterval   time.Duration
	eurekaPath        string
}

func New(eurekaUrls []string, opts ...Option) (*Registry, error) {
	r := &Registry{
		ctx:               context.Background(),
		heartbeatInterval: heartbeatTime,
		refreshInterval:   refreshTime,
		eurekaPath:        "eureka/v2",
	}

	for _, o := range opts {
		o(r)
	}

	client := NewClient(eurekaUrls, WithHeartbeatInterval(r.heartbeatInterval), WithClientContext(r.ctx), WithNamespace(r.eurekaPath))
	r.api = NewAPI(r.ctx, client, r.refreshInterval)
	return r, nil
}

// 这里的Context是每个注册器独享的
func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	return r.api.Register(ctx, service.Name, r.Endpoints(service)...)
}

// Deregister registry service to zookeeper.
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return r.api.Deregister(ctx, r.Endpoints(service))
}

// GetService get services from zookeeper
func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	instances := r.api.GetService(ctx, serviceName)
	items := make([]*registry.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		items = append(items, &registry.ServiceInstance{
			ID:        instance.Metadata["ID"],
			Name:      instance.Metadata["Name"],
			Version:   instance.Metadata["Version"],
			Endpoints: []string{instance.Metadata["Endpoints"]},
			Metadata:  instance.Metadata,
		})
	}

	return items, nil
}

// watch 是独立的ctx
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatch(ctx, r.api, serviceName)
}

func (r *Registry) Endpoints(service *registry.ServiceInstance) []Endpoint {
	var (
		res   = []Endpoint{}
		start int
	)
	for _, ep := range service.Endpoints {
		start = strings.Index(ep, "//")
		end := strings.LastIndex(ep, ":")
		appID := strings.ToUpper(service.Name)
		ip := ep[start+2 : end]
		sport := ep[end+1:]
		port, _ := strconv.Atoi(sport)
		securePort := 443
		homePageURL := fmt.Sprintf("%s/", ep)
		statusPageURL := fmt.Sprintf("%s/info", ep)
		healthCheckURL := fmt.Sprintf("%s/health", ep)
		instanceID := strings.Join([]string{ip, appID, sport}, ":")
		metadata := make(map[string]string)
		if len(service.Metadata) > 0 {
			metadata = service.Metadata
		}
		if s, ok := service.Metadata["securePort"]; ok {
			securePort, _ = strconv.Atoi(s)
		}
		if s, ok := service.Metadata["homePageURL"]; ok {
			homePageURL = s
		}
		if s, ok := service.Metadata["statusPageURL"]; ok {
			statusPageURL = s
		}
		if s, ok := service.Metadata["healthCheckURL"]; ok {
			healthCheckURL = s
		}
		metadata["ID"] = service.ID
		metadata["Name"] = service.Name
		metadata["Version"] = service.Version
		metadata["Endpoints"] = ep
		metadata["agent"] = "go-eureka-client"
		res = append(res, Endpoint{
			AppID:          appID,
			IP:             ip,
			Port:           port,
			SecurePort:     securePort,
			HomePageURL:    homePageURL,
			StatusPageURL:  statusPageURL,
			HealthCheckURL: healthCheckURL,
			InstanceID:     instanceID,
			MetaData:       metadata,
		})
	}

	return res
}
