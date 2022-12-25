package polaris

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
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

// _instanceIDSeparator . Instance id Separator.
const _instanceIDSeparator = "-"

type registryOptions struct {
	// required, namespace in polaris
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

// WithRegistryNamespace with Namespace option.
func WithRegistryNamespace(namespace string) RegistryOption {
	return func(o *registryOptions) { o.Namespace = namespace }
}

// WithServiceToken with ServiceToken option.
func WithServiceToken(serviceToken string) RegistryOption {
	return func(o *registryOptions) { o.ServiceToken = serviceToken }
}

// WithWeight with Weight option.
func WithWeight(weight int) RegistryOption {
	return func(o *registryOptions) { o.Weight = weight }
}

// WithHealthy with Healthy option.
func WithHealthy(healthy bool) RegistryOption {
	return func(o *registryOptions) { o.Healthy = healthy }
}

// WithIsolate with Isolate option.
func WithIsolate(isolate bool) RegistryOption {
	return func(o *registryOptions) { o.Isolate = isolate }
}

// WithTTL with TTL option.
func WithTTL(TTL int) RegistryOption {
	return func(o *registryOptions) { o.TTL = TTL }
}

// WithTimeout with Timeout option.
func WithTimeout(timeout time.Duration) RegistryOption {
	return func(o *registryOptions) { o.Timeout = timeout }
}

// WithRetryCount with RetryCount option.
func WithRetryCount(retryCount int) RegistryOption {
	return func(o *registryOptions) { o.RetryCount = retryCount }
}

func (p *Polaris) Registry(opts ...RegistryOption) (r *Registry) {
	op := registryOptions{
		Namespace:    "default",
		ServiceToken: "",
		Weight:       0,
		Priority:     0,
		Healthy:      true,
		Isolate:      false,
		TTL:          0,
		Timeout:      0,
		RetryCount:   0,
	}
	for _, option := range opts {
		option(&op)
	}
	return &Registry{
		opt:      op,
		provider: p.registry,
		consumer: p.discovery,
	}
}

// Register the registration.
func (r *Registry) Register(_ context.Context, instance *registry.ServiceInstance) error {
	ids := make([]string, 0, len(instance.Endpoints))
	if instance.ID == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		instance.ID = id.String()
	}

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
		instance.Metadata["id"] = instance.ID
		if _, ok := instance.Metadata["weight"]; !ok {
			instance.Metadata["weight"] = strconv.Itoa(r.opt.Weight)
		}
		weight, _ := strconv.Atoi(instance.Metadata["weight"])

		m, err := r.provider.RegisterInstance(
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
					Metadata:     instance.Metadata,
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
		ids = append(ids, m.InstanceID)
	}
	// need to set InstanceID for Deregister
	instance.ID = strings.Join(ids, _instanceIDSeparator)
	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(_ context.Context, serviceInstance *registry.ServiceInstance) error {
	split := strings.Split(serviceInstance.ID, _instanceIDSeparator)
	for i, endpoint := range serviceInstance.Endpoints {
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
					Service:      serviceInstance.Name + u.Scheme,
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					InstanceID:   split[i],
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
			Service:    serviceName,
			Namespace:  r.opt.Namespace,
			Timeout:    &r.opt.Timeout,
			RetryCount: &r.opt.RetryCount,
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
		if v, ok := m[instance.GetHost()+instance.GetMetadata()["id"]]; ok {
			m[instance.GetHost()+instance.GetMetadata()["id"]] = append(v, instance)
		} else {
			m[instance.GetHost()+instance.GetMetadata()["id"]] = []model.Instance{instance}
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
						if v, ok := w.ServiceInstances[instance.GetHost()+instance.GetMetadata()["id"]]; ok {
							var nv []model.Instance
							for _, m := range v {
								if m.GetId() != instance.GetId() {
									nv = append(nv, instance)
								}
							}
							for index, ins := range w.service.Instances {
								if instance.GetId() == ins.GetId() {
									// remove equal
									if len(w.service.Instances) <= 1 {
										w.service.Instances = w.service.Instances[0:0]
										continue
									}
									w.service.Instances = append(w.service.Instances[:index], w.service.Instances[index+1:]...)
								}
							}
							w.ServiceInstances[instance.GetHost()+instance.GetMetadata()["id"]] = nv
						}
					}
				}
				// handle UpdateEvent
				if instanceEvent.UpdateEvent != nil {
					for _, update := range instanceEvent.UpdateEvent.UpdateList {
						if v, ok := w.ServiceInstances[update.After.GetHost()+update.After.GetMetadata()["id"]]; ok {
							nv := []model.Instance{update.After}
							for _, m := range v {
								if m.GetId() != update.After.GetId() {
									// Insert directly those not updated this time
									nv = append(nv, m)
								}
							}
							w.ServiceInstances[update.After.GetHost()+update.After.GetMetadata()["id"]] = nv
						} else {
							w.ServiceInstances[update.After.GetHost()+update.After.GetMetadata()["id"]] = []model.Instance{update.After}
						}
						for i, serviceInstance := range w.service.Instances {
							for _, update := range instanceEvent.UpdateEvent.UpdateList {
								if serviceInstance.GetId() == update.Before.GetId() {
									w.service.Instances[i] = update.After
								}
							}
						}
					}
				}
				// handle AddEvent
				if instanceEvent.AddEvent != nil {
					for _, instance := range instanceEvent.AddEvent.Instances {
						if v, ok := w.ServiceInstances[instance.GetHost()+instance.GetMetadata()["id"]]; ok {
							var nv []model.Instance
							m := map[string]model.Instance{}
							for _, ins := range v {
								m[ins.GetId()] = ins
							}
							m[instance.GetId()] = instance
							for _, ins := range m {
								nv = append(nv, ins)
							}
							w.ServiceInstances[instance.GetHost()+instance.GetMetadata()["id"]] = nv
						} else {
							w.ServiceInstances[instance.GetHost()+instance.GetMetadata()["id"]] = []model.Instance{instance}
						}
						w.service.Instances = append(w.service.Instances, instanceEvent.AddEvent.Instances...)
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
			ID:        inss[0].GetMetadata()["id"],
			Name:      inss[0].GetService(),
			Version:   inss[0].GetVersion(),
			Metadata:  inss[0].GetMetadata(),
			Endpoints: []string{fmt.Sprintf("%s://%s:%d", inss[0].GetProtocol(), inss[0].GetHost(), inss[0].GetPort())},
		}
		for i := 1; i < len(inss); i++ {
			ins.Endpoints = append(ins.Endpoints, fmt.Sprintf("%s://%s:%d", inss[i].GetProtocol(), inss[i].GetHost(), inss[i].GetPort()))
		}
		serviceInstances = append(serviceInstances, ins)
	}
	return serviceInstances
}
