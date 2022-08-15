package grpc

import (
	"sync"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/transport"

	gBalancer "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

const (
	globalSelectorName = "global"
)

var (
	_ base.PickerBuilder = &Builder{}
	_ gBalancer.Picker   = &Picker{}

	mu sync.Mutex
)

func init() {
	// inject global grpc balancer
	SetGlobalBalancer(globalSelectorName, transport.GlobalSelector())
}

// SetGlobalBalancer set grpc balancer with scheme.
func SetGlobalBalancer(scheme string, builder selector.Builder) {
	mu.Lock()
	defer mu.Unlock()

	b := base.NewBalancerBuilder(
		scheme,
		&Builder{builder: builder},
		base.Config{HealthCheck: true},
	)
	gBalancer.Register(b)
}

// Builder is grpc balancer builder.
type Builder struct {
	builder selector.Builder
}

// Build creates a grpc Picker.
func (b *Builder) Build(info base.PickerBuildInfo) gBalancer.Picker {
	if len(info.ReadySCs) == 0 {
		// Block the RPC until a new picker is available via UpdateState().
		return base.NewErrPicker(gBalancer.ErrNoSubConnAvailable)
	}

	nodes := make([]selector.Node, 0)
	for conn, info := range info.ReadySCs {
		ins, _ := info.Address.Attributes.Value("rawServiceInstance").(*registry.ServiceInstance)
		nodes = append(nodes, &grpcNode{
			Node:    selector.NewNode("grpc", info.Address.Addr, ins),
			subConn: conn,
		})
	}
	p := &Picker{
		selector: b.builder.Build(),
	}
	p.selector.Apply(nodes)
	return p
}

// Picker is a grpc picker.
type Picker struct {
	selector selector.Selector
}

// Pick pick instances.
func (p *Picker) Pick(info gBalancer.PickInfo) (gBalancer.PickResult, error) {
	var filters []selector.NodeFilter
	if tr, ok := transport.FromClientContext(info.Ctx); ok {
		if gtr, ok := tr.(*Transport); ok {
			filters = gtr.SelectFilters()
		}
	}

	n, done, err := p.selector.Select(info.Ctx, selector.WithFilter(filters...))
	if err != nil {
		return gBalancer.PickResult{}, err
	}

	return gBalancer.PickResult{
		SubConn: n.(*grpcNode).subConn,
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

type grpcNode struct {
	selector.Node
	subConn gBalancer.SubConn
}
