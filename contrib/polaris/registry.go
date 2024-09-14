package polaris

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/go-kratos/kratos/v2/registry"
)

var (
	_ registry.Registrar = (*Registry)(nil)
	_ registry.Discovery = (*Registry)(nil)
)

type registryOptions struct {
	// required, testNamespace in polaris
	Namespace string

	// required, service access token
	ServiceToken string

	// service weight in polaris. Default value is 100, 0 <= weight <= 10000
	Weight int

	// service priority. Default value is 0. The smaller the value, the lower the priority
	Priority int

	// To show service is healthy or not. Default value is True .
	Healthy bool

	// To show service is isolate or not. Default value is False .
	Isolate bool

	// TTL timeout. if node needs to use heartbeat to report,required. If not set,server will throw ErrorCode-400141
	TTL int

	// optional, Timeout for single query. Default value is global config
	// Total is (1+RetryCount) * Timeout
	Timeout time.Duration

	// optional, retry count. Default value is global config
	RetryCount int
}

// RegistryOption is polaris option.
type RegistryOption func(o *registryOptions)

// Registry is polaris registry.
type Registry struct {
	opt      registryOptions
	provider polaris.ProviderAPI
	consumer polaris.ConsumerAPI
}

// WithRegistryServiceToken with ServiceToken option.
func WithRegistryServiceToken(serviceToken string) RegistryOption {
	return func(o *registryOptions) { o.ServiceToken = serviceToken }
}

// WithRegistryWeight with Weight option.
func WithRegistryWeight(weight int) RegistryOption {
	return func(o *registryOptions) { o.Weight = weight }
}

// WithRegistryHealthy with Healthy option.
func WithRegistryHealthy(healthy bool) RegistryOption {
	return func(o *registryOptions) { o.Healthy = healthy }
}

// WithRegistryIsolate with Isolate option.
func WithRegistryIsolate(isolate bool) RegistryOption {
	return func(o *registryOptions) { o.Isolate = isolate }
}

// WithRegistryTTL with TTL option.
func WithRegistryTTL(TTL int) RegistryOption {
	return func(o *registryOptions) { o.TTL = TTL }
}

// WithRegistryTimeout with Timeout option.
func WithRegistryTimeout(timeout time.Duration) RegistryOption {
	return func(o *registryOptions) { o.Timeout = timeout }
}

// WithRegistryRetryCount with RetryCount option.
func WithRegistryRetryCount(retryCount int) RegistryOption {
	return func(o *registryOptions) { o.RetryCount = retryCount }
}

// Register the registration.
func (r *Registry) Register(_ context.Context, instance *registry.ServiceInstance) error {
	id := uuid.NewString()
	for _, endpoint := range instance.Endpoints {
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}

		portNum, err := strconv.Atoi(port)
		if err != nil {
			return err
		}

		// metadata
		rmd := mapClone(instance.Metadata)
		if rmd == nil {
			rmd = make(map[string]string)
		}
		rmd["merge"] = id
		if _, ok := rmd["weight"]; !ok {
			rmd["weight"] = strconv.Itoa(r.opt.Weight)
		}

		weight, _ := strconv.Atoi(rmd["weight"])

		_, err = r.provider.RegisterInstance(
			&polaris.InstanceRegisterRequest{
				InstanceRegisterRequest: model.InstanceRegisterRequest{
					Service:      instance.Name,
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					Host:         host,
					Port:         portNum,
					Protocol:     &u.Scheme,
					Weight:       &weight,
					Priority:     &r.opt.Priority,
					Version:      &instance.Version,
					Metadata:     rmd,
					Healthy:      &r.opt.Healthy,
					Isolate:      &r.opt.Isolate,
					TTL:          &r.opt.TTL,
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(_ context.Context, serviceInstance *registry.ServiceInstance) error {
	for _, endpoint := range serviceInstance.Endpoints {
		// get url
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		// get host and port
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}

		// port to int
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		// Deregister
		err = r.provider.Deregister(
			&polaris.InstanceDeRegisterRequest{
				InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
					Service:      serviceInstance.Name,
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					Host:         host,
					Port:         portNum,
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetService return the service instances in memory according to the service name.
func (r *Registry) GetService(_ context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	// get all instances
	instancesResponse, err := r.consumer.GetInstances(&polaris.GetInstancesRequest{
		GetInstancesRequest: model.GetInstancesRequest{
			Service:         serviceName,
			Namespace:       r.opt.Namespace,
			Timeout:         &r.opt.Timeout,
			RetryCount:      &r.opt.RetryCount,
			SkipRouteFilter: true,
		},
	})
	if err != nil {
		return nil, err
	}

	serviceInstances := instancesToServiceInstances(merge(instancesResponse.GetInstances()))

	return serviceInstances, nil
}

func merge(instances []model.Instance) map[string][]model.Instance {
	m := make(map[string][]model.Instance)
	for _, instance := range instances {
		if v, ok := m[instance.GetMetadata()["merge"]]; ok {
			m[instance.GetMetadata()["merge"]] = append(v, instance)
		} else {
			m[instance.GetMetadata()["merge"]] = []model.Instance{instance}
		}
	}
	return m
}

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatcher(ctx, r.opt.Namespace, serviceName, r.consumer)
}

type Watcher struct {
	ServiceName      string
	Namespace        string
	Ctx              context.Context
	Cancel           context.CancelFunc
	Channel          <-chan model.SubScribeEvent
	service          *model.InstancesResponse
	ServiceInstances map[string][]model.Instance
	first            bool
}

func newWatcher(ctx context.Context, namespace string, serviceName string, consumer polaris.ConsumerAPI) (*Watcher, error) {
	watchServiceResponse, err := consumer.WatchService(&polaris.WatchServiceRequest{
		WatchServiceRequest: model.WatchServiceRequest{
			Key: model.ServiceKey{
				Namespace: namespace,
				Service:   serviceName,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		Namespace:        namespace,
		ServiceName:      serviceName,
		Channel:          watchServiceResponse.EventChannel,
		service:          watchServiceResponse.GetAllInstancesResp,
		ServiceInstances: merge(watchServiceResponse.GetAllInstancesResp.GetInstances()),
	}
	w.Ctx, w.Cancel = context.WithCancel(ctx)
	return w, nil
}

// Next returns services in the following two cases:
// 1.the first time to watch and the service instance list is not empty.
// 2.any service instance changes found.
// if the above two conditions are not met, it will block until context deadline exceeded or canceled
func (w *Watcher) Next() ([]*registry.ServiceInstance, error) {
	if !w.first {
		w.first = true
		if len(w.ServiceInstances) > 0 {
			return instancesToServiceInstances(w.ServiceInstances), nil
		}
	}
	select {
	case <-w.Ctx.Done():
		return nil, w.Ctx.Err()
	case event := <-w.Channel:
		if event.GetSubScribeEventType() == model.EventInstance {
			// this always true, but we need to check it to make sure EventType not change
			if instanceEvent, ok := event.(*model.InstanceEvent); ok {
				// handle DeleteEvent
				if instanceEvent.DeleteEvent != nil {
					for _, instance := range instanceEvent.DeleteEvent.Instances {
						delete(w.ServiceInstances, instance.GetMetadata()["merge"])
					}
				}
				// handle UpdateEvent
				if instanceEvent.UpdateEvent != nil {
					for _, update := range instanceEvent.UpdateEvent.UpdateList {
						if v, ok := w.ServiceInstances[update.After.GetMetadata()["merge"]]; ok {
							var nv []model.Instance
							m := map[string]model.Instance{}
							for _, ins := range v {
								m[ins.GetId()] = ins
							}
							m[update.After.GetId()] = update.After
							for _, ins := range m {
								if ins.IsHealthy() {
									nv = append(nv, ins)
								}
							}
							w.ServiceInstances[update.After.GetMetadata()["merge"]] = nv
							if len(nv) == 0 {
								delete(w.ServiceInstances, update.After.GetMetadata()["merge"])
							}
						} else {
							if update.After.IsHealthy() {
								w.ServiceInstances[update.After.GetMetadata()["merge"]] = []model.Instance{update.After}
							}
						}
					}
				}
				// handle AddEvent
				if instanceEvent.AddEvent != nil {
					for _, instance := range instanceEvent.AddEvent.Instances {
						if v, ok := w.ServiceInstances[instance.GetMetadata()["merge"]]; ok {
							var nv []model.Instance
							m := map[string]model.Instance{}
							for _, ins := range v {
								m[ins.GetId()] = ins
							}
							m[instance.GetId()] = instance
							for _, ins := range m {
								if ins.IsHealthy() {
									nv = append(nv, ins)
								}
							}
							if len(nv) != 0 {
								w.ServiceInstances[instance.GetMetadata()["merge"]] = nv
							}
						} else {
							if instance.IsHealthy() {
								w.ServiceInstances[instance.GetMetadata()["merge"]] = []model.Instance{instance}
							}
						}
					}
				}
			}
			return instancesToServiceInstances(w.ServiceInstances), nil
		}
	}
	return instancesToServiceInstances(w.ServiceInstances), nil
}

// Stop close the watcher.
func (w *Watcher) Stop() error {
	w.Cancel()
	return nil
}

func instancesToServiceInstances(instances map[string][]model.Instance) []*registry.ServiceInstance {
	serviceInstances := make([]*registry.ServiceInstance, 0, len(instances))
	for _, inss := range instances {
		if len(inss) == 0 {
			continue
		}
		ins := &registry.ServiceInstance{
			ID:       inss[0].GetId(),
			Name:     inss[0].GetService(),
			Version:  inss[0].GetVersion(),
			Metadata: inss[0].GetMetadata(),
		}
		for _, item := range inss {
			if item.IsHealthy() {
				ins.Endpoints = append(ins.Endpoints, fmt.Sprintf("%s://%s:%d", item.GetProtocol(), item.GetHost(), item.GetPort()))
			}
		}
		if len(ins.Endpoints) != 0 {
			serviceInstances = append(serviceInstances, ins)
		}
	}
	return serviceInstances
}

// Clone returns a copy of m. This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func mapClone[M ~map[K]V, K comparable, V any](m M) M {
	// Preserve nil in case it matters.
	if m == nil {
		return nil
	}
	// Make a shallow copy of the map.
	m2 := make(M, len(m))
	for k, v := range m {
		m2[k] = v
	}
	return m2
}
