package polaris

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"
)

var _ registry.Registrar = (*Registry)(nil)

// _instanceIDSeparator . Instance id Separator.
const _instanceIDSeparator = "-"

type options struct {

	// required, namespace in polaris
	Namespace string

	// required, service access token
	ServiceToken string

	// optional, protocol in polaris. Default value is nil, it means use protocol config in service
	Protocol *string

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

// Option is polaris option.
type Option func(o *options)

// Registry is polaris registry.
type Registry struct {
	opt      options
	provider api.ProviderAPI
}

// WithNamespace with default Namespace option.
func WithNamespace(namespace string) Option {
	return func(o *options) { o.Namespace = namespace }
}

// WithServiceToken with default ServiceToken option.
func WithServiceToken(serviceToken string) Option {
	return func(o *options) { o.ServiceToken = serviceToken }
}

// WithProtocol with default Protocol option.
func WithProtocol(protocol string) Option {
	return func(o *options) { o.Protocol = &protocol }
}

// WithDefaultWeight with default Weight option.
func WithDefaultWeight(weight int) Option {
	return func(o *options) { o.Weight = weight }
}

// WithHealthy with default Healthy option.
func WithHealthy(healthy bool) Option {
	return func(o *options) { o.Healthy = healthy }
}

// WithIsolate with default Isolate option.
func WithIsolate(isolate bool) Option {
	return func(o *options) { o.Isolate = isolate }
}

// WithDefaultTTL with default TTL option.
func WithDefaultTTL(TTL int) Option {
	return func(o *options) { o.TTL = TTL }
}

// WithTimeout with default Timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.Timeout = timeout }
}

// WithRetryCount with default RetryCount option.
func WithRetryCount(retryCount int) Option {
	return func(o *options) { o.RetryCount = retryCount }
}

func NewRegistry(provider api.ProviderAPI, opts ...Option) (r *Registry) {
	op := options{
		Namespace:    "default",
		ServiceToken: "",
		Protocol:     nil,
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
		provider: provider,
	}
}

// Register the registration.
func (r *Registry) Register(_ context.Context, serviceInstance *registry.ServiceInstance) error {
	ids := make([]string, 0, len(serviceInstance.Endpoints))
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

		// medata
		var rmd map[string]string
		if serviceInstance.Metadata == nil {
			rmd = map[string]string{
				"kind":    u.Scheme,
				"version": serviceInstance.Version,
			}
		} else {
			rmd = make(map[string]string, len(serviceInstance.Metadata)+2)
			for k, v := range serviceInstance.Metadata {
				rmd[k] = v
			}
			rmd["kind"] = u.Scheme
			rmd["version"] = serviceInstance.Version
		}
		// Register
		fmt.Println(serviceInstance.Name + u.Scheme)
		service, err := r.provider.Register(
			&api.InstanceRegisterRequest{
				InstanceRegisterRequest: model.InstanceRegisterRequest{
					Service:      serviceInstance.Name + u.Scheme,
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					Host:         host,
					Port:         portNum,
					Protocol:     r.opt.Protocol,
					Weight:       &r.opt.Weight,
					Priority:     &r.opt.Priority,
					Version:      &serviceInstance.Version,
					Metadata:     rmd,
					Healthy:      &r.opt.Healthy,
					Isolate:      &r.opt.Isolate,
					TTL:          &r.opt.TTL,
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			})
		if err != nil {
			return err
		}
		ids = append(ids, service.InstanceID)
	}
	// need to set InstanceID for Deregister
	serviceInstance.ID = strings.Join(ids, _instanceIDSeparator)
	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, serviceInstance *registry.ServiceInstance) error {
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
			&api.InstanceDeRegisterRequest{
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
