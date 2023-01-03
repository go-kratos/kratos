package polaris

import (
	"errors"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/api"
)

type Polaris struct {
	router    polaris.RouterAPI
	config    polaris.ConfigAPI
	limit     polaris.LimitAPI
	registry  polaris.ProviderAPI
	discovery polaris.ConsumerAPI
	namespace string
}

// Option is polaris option.
type Option func(o *Polaris)

// WithNamespace with polaris global testNamespace
func WithNamespace(ns string) Option {
	return func(o *Polaris) {
		o.namespace = ns
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
	}
	for _, option := range opts {
		option(&op)
	}
	return op
}

func (p *Polaris) Config(opts ...ConfigOption) (config.Source, error) {
	options := &configOptions{
		namespace: p.namespace,
		fileGroup: "",
		fileName:  "",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.fileGroup == "" {
		return nil, errors.New("testFileGroup invalid")
	}

	if options.fileName == "" {
		return nil, errors.New("fileName invalid")
	}

	return &source{
		client:  p.config,
		options: options,
	}, nil
}

func (p *Polaris) Registry(opts ...RegistryOption) (r *Registry) {
	op := registryOptions{
		Namespace:    p.namespace,
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
