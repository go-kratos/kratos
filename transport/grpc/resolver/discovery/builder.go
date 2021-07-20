package discovery

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"google.golang.org/grpc/resolver"
)

const name = "discovery"

// Option is builder option.
type Option func(o *builder)

// WithLogger with builder logger.
func WithLogger(logger log.Logger) Option {
	return func(b *builder) {
		b.logger = logger
	}
}

// WithTimeout with timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(b *builder) {
		b.timeout = timeout
	}
}

type builder struct {
	discoverer registry.Discovery
	logger     log.Logger
	timeout    time.Duration
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(d registry.Discovery, opts ...Option) resolver.Builder {
	b := &builder{
		discoverer: d,
		logger:     log.DefaultLogger,
		timeout:    time.Second * 10,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()
	w, err := b.discoverer.Watch(ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}

	r := &discoveryResolver{
		w:      w,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
		log:    log.NewHelper(b.logger),
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	go r.watch()

	return r, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
	return name
}
