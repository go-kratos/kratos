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

// InstanceIDSeparator . Instance Id Separator.
const InstanceIDSeparator = "-"

type options struct {
	// 必选，服务名
	NameSpace string

	// 必选，服务访问Token
	ServiceToken string

	// 以下字段可选，默认nil表示客户端不配置，使用服务端配置
	// 服务协议
	Protocol *string

	// 服务权重，默认100，范围0-10000
	Weight int

	// 实例优先级，默认为0，数值越小，优先级越高
	Priority int

	// 该服务实例是否健康，默认健康
	Healthy bool

	// 该服务实例是否隔离，默认不隔离
	Isolate bool

	// ttl超时时间，如果节点要调用heartbeat上报，则必须填写，否则会400141错误码，单位：秒
	TTL int

	// 可选，单次查询超时时间，默认直接获取全局的超时配置
	// 用户总最大超时时间为(1+RetryCount) * Timeout
	Timeout time.Duration

	// 可选，重试次数，默认直接获取全局的超时配置
	RetryCount int
}

// Option is polaris option.
type Option func(o *options)

// Registry is polaris registry.
type Registry struct {
	opt      options
	provider api.ProviderAPI
}

// WithDefaultNamespace with default NameSpace option.
func WithDefaultNamespace(nameSpace string) Option {
	return func(o *options) { o.NameSpace = nameSpace }
}

// WithDefaultServiceToken with default ServiceToken option.
func WithDefaultServiceToken(serviceToken string) Option {
	return func(o *options) { o.ServiceToken = serviceToken }
}

// WithDefaultProtocol with default Protocol option.
func WithDefaultProtocol(protocol string) Option {
	return func(o *options) { o.Protocol = &protocol }
}

// WithDefaultWeight with default Weight option.
func WithDefaultWeight(weight int) Option {
	return func(o *options) { o.Weight = weight }
}

// WithDefaultHealthy with default Healthy option.
func WithDefaultHealthy(healthy bool) Option {
	return func(o *options) { o.Healthy = healthy }
}

// WithDefaultIsolate with default Isolate option.
func WithDefaultIsolate(isolate bool) Option {
	return func(o *options) { o.Isolate = isolate }
}

// WithDefaultTTL with default TTL option.
func WithDefaultTTL(TTL int) Option {
	return func(o *options) { o.TTL = TTL }
}

// WithDefaultTimeout with default Timeout option.
func WithDefaultTimeout(timeout time.Duration) Option {
	return func(o *options) { o.Timeout = timeout }
}

// WithDefaultRetryCount with default RetryCount option.
func WithDefaultRetryCount(retryCount int) Option {
	return func(o *options) { o.RetryCount = retryCount }
}

func New(provider api.ProviderAPI, opts ...Option) (r *Registry) {
	op := options{
		NameSpace:    "default",
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
					Namespace:    r.opt.NameSpace,
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
	serviceInstance.ID = strings.Join(ids, InstanceIDSeparator)
	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, serviceInstance *registry.ServiceInstance) error {
	split := strings.Split(serviceInstance.ID, InstanceIDSeparator)
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
					Namespace:    r.opt.NameSpace,
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
