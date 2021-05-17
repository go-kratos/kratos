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
	discoverer registry.Discovery
	logger     log.Logger
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
	b := &builder{
		discoverer: d,
		logger:     log.DefaultLogger,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (d *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	w, err := d.discoverer.Watch(context.Background(), target.Endpoint)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := &discoveryResolver{
		w:      w,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
		log:    log.NewHelper(d.logger),
	}
	go r.watch()
	return r, nil
}

func (d *builder) Scheme() string {
	return name
}
