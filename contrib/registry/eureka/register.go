package eureka

import (
	"context"
	"strings"

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

func WithHeartbeat(interval string) Option {
	return func(o *Registry) { o.heartbeatInterval = interval }
}

func WithRefresh(interval string) Option {
	return func(o *Registry) { o.refreshInterval = interval }
}

type Registry struct {
	ctx               context.Context
	api               *eurekaApi
	heartbeatInterval string
	refreshInterval   string
}

func New(eurekaUrls []string, opts ...Option) (*Registry, error) {
	r := &Registry{
		ctx:               context.Background(),
		heartbeatInterval: "10s",
		refreshInterval:   "30s",
	}

	for _, o := range opts {
		o(r)
	}

	client := NewEurekaClient(eurekaUrls, WithHeartbeatInterval(r.heartbeatInterval), WithCtx(r.ctx))
	r.api = NewEurekaApi(client, r.refreshInterval)
	return r, nil
}

func (r *Registry) buildInstance(list []Instance) map[string][]*registry.ServiceInstance {
	items := make(map[string][]*registry.ServiceInstance)

	for _, instance := range list {
		item := &registry.ServiceInstance{
			ID:        instance.Metadata["ID"],
			Name:      instance.Metadata["Name"],
			Version:   instance.Metadata["Version"],
			Metadata:  instance.Metadata,
			Endpoints: []string{instance.Metadata["Endpoints"]},
		}
		items[instance.App] = append(items[instance.App], item)
	}

	return items
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
	var res = []Endpoint{}
	for _, ep := range service.Endpoints {
		var start int
		start = strings.Index(ep, "//")
		end := strings.LastIndex(ep, ":")
		appID := strings.ToUpper(service.Name)
		ip := ep[start+2 : end]
		port := ep[end+1:]
		securePort := "443"
		homePageUrl := "/"
		statusPageUrl := "/info"
		healthCheckUrl := "/health"
		instanceId := strings.Join([]string{ip, appID, port}, ":")
		metadata := make(map[string]string)
		if service.Metadata != nil {
			metadata = service.Metadata
		}
		if s, ok := service.Metadata["securePort"]; ok {
			securePort = s
		}
		if s, ok := service.Metadata["homePageUrl"]; ok {
			homePageUrl = s
		}
		if s, ok := service.Metadata["statusPageUrl"]; ok {
			statusPageUrl = s
		}
		if s, ok := service.Metadata["healthCheckUrl"]; ok {
			healthCheckUrl = s
		}
		metadata["ID"] = service.ID
		metadata["Name"] = service.Name
		metadata["Version"] = service.Version
		metadata["Endpoints"] = ep
		res = append(res, Endpoint{
			AppID:          appID,
			IP:             ip,
			Port:           port,
			SecurePort:     securePort,
			HomePageUrl:    homePageUrl,
			StatusPageUrl:  statusPageUrl,
			HealthCheckUrl: healthCheckUrl,
			InstanceID:     instanceId,
			MetaData:       metadata,
		})
	}

	return res
}
