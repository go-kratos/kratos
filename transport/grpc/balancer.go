package grpc

import (
	"sync"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport"

	gBalancer "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

var (
	_ base.PickerBuilder = &Builder{}
	_ gBalancer.Picker   = &Picker{}

	mu sync.Mutex
)

func init() {
	// inject global grpc balancer
	SetGlobalBalancer(random.Name, random.New())
	SetGlobalBalancer(wrr.Name, wrr.New())
	SetGlobalBalancer(p2c.Name, p2c.New())
}

// SetGlobalBalancer set grpc balancer with scheme.
func SetGlobalBalancer(scheme string, selector selector.Selector) {
	mu.Lock()
	defer mu.Unlock()

	b := base.NewBalancerBuilder(
		scheme,
		&Builder{selector},
		base.Config{HealthCheck: true},
	)
	gBalancer.Register(b)
}

// Builder is grpc balancer builder.
type Builder struct {
	selector selector.Selector
}

// Build creates a grpc Picker.
func (b *Builder) Build(info base.PickerBuildInfo) gBalancer.Picker {
	nodes := make([]selector.Node, 0)
	subConns := make(map[string]gBalancer.SubConn)
	for conn, info := range info.ReadySCs {
		if _, ok := subConns[info.Address.Addr]; ok {
			continue
		}
		subConns[info.Address.Addr] = conn

		ins, _ := info.Address.Attributes.Value("rawServiceInstance").(*registry.ServiceInstance)
		nodes = append(nodes, node.New(info.Address.Addr, ins))
	}
	p := &Picker{
		selector: b.selector,
		subConns: subConns,
	}
	p.selector.Apply(nodes)
	return p
}

// Picker is a grpc picker.
type Picker struct {
	subConns map[string]gBalancer.SubConn
	selector selector.Selector
}

// Pick pick instances.
func (p *Picker) Pick(info gBalancer.PickInfo) (gBalancer.PickResult, error) {
	var filters []selector.Filter
	if tr, ok := transport.FromClientContext(info.Ctx); ok {
		if gtr, ok := tr.(*Transport); ok {
			filters = gtr.Filters()
		}
	}

	n, done, err := p.selector.Select(info.Ctx, selector.WithFilter(filters...))
	if err != nil {
		return gBalancer.PickResult{}, err
	}
	sub := p.subConns[n.Address()]

	return gBalancer.PickResult{
		SubConn: sub,
		Done: func(di gBalancer.DoneInfo) {
			done(info.Ctx, selector.DoneInfo{
				Err:           di.Err,
				BytesSent:     di.BytesSent,
				BytesReceived: di.BytesReceived,
				ReplyMeta:     Trailer(di.Trailer),
			})
		},
	}, nil
}

// Trailer is a grpc trailder MD.
type Trailer metadata.MD

// Get get a grpc trailer value.
func (t Trailer) Get(k string) string {
	v := metadata.MD(t).Get(k)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}
