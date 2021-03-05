package discovery

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"google.golang.org/grpc/resolver"
)

const name = "discovery"

// Option is builder option.
type Option func(o *builder)

// WithLogger with builder logger.
func WithLogger(logger log.Logger) Option {
	return func(o *builder) {
		o.logger = logger
	}
}

type builder struct {
	instancer registry.Instancer
	logger    log.Logger
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(in registry.Instancer, opts ...Option) resolver.Builder {
	b := &builder{
		instancer: in,
		logger:    log.DefaultLogger,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (d *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	w, err := d.instancer.Watch(context.Background(), target.Endpoint)
	if err != nil {
		return nil, err
	}
	r := &discoveryResolver{
		w:   w,
		cc:  cc,
		log: log.NewHelper("grpc/resolver/discovery", d.logger),
	}
	go r.watch()
	return r, nil
}

func (d *builder) Scheme() string {
	return name
}
