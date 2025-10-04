package polaris

import (
	"errors"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/api"

	"github.com/go-kratos/kratos/v2/config"
)

type Polaris struct {
	router    polaris.RouterAPI
	config    polaris.ConfigAPI
	limit     polaris.LimitAPI
	registry  polaris.ProviderAPI
	discovery polaris.ConsumerAPI
	namespace string
	service   string
}

// Option is polaris option.
type Option func(o *Polaris)

// WithNamespace with polaris global testNamespace
func WithNamespace(ns string) Option {
	return func(o *Polaris) {
		o.namespace = ns
	}
}

// WithService set the current service name
func WithService(service string) Option {
	return func(o *Polaris) {
		o.service = service
	}
}

// New polaris Service governance.
func New(sdk api.SDKContext, opts ...Option) Polaris {
	op := Polaris{
		router:    polaris.NewRouterAPIByContext(sdk),
		config:    polaris.NewConfigAPIByContext(sdk),
		limit:     polaris.NewLimitAPIByContext(sdk),
		registry:  polaris.NewProviderAPIByContext(sdk),
		discovery: polaris.NewConsumerAPIByContext(sdk),
		namespace: "default",
	}
	for _, option := range opts {
		option(&op)
	}
	return op
}

func (p *Polaris) Config(opts ...ConfigOption) (config.Source, error) {
	options := &configOptions{
		namespace: p.namespace,
	}

	for _, opt := range opts {
		opt(options)
	}

	if len(options.files) == 0 {
		return nil, errors.New("fileNames invalid")
	}

	return &source{
		client:  p.config,
		options: options,
	}, nil
}

func (p *Polaris) Registry(opts ...RegistryOption) (r *Registry) {
	op := registryOptions{
		Namespace: p.namespace,
		Healthy:   true,
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

func (p *Polaris) Limiter(opts ...LimiterOption) (r Limiter) {
	op := limiterOptions{
		namespace: p.namespace,
		service:   p.service,
	}
	for _, option := range opts {
		option(&op)
	}
	return Limiter{
		limitAPI: p.limit,
		opts:     op,
	}
}
