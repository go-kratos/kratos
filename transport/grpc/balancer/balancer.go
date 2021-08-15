package multi

import (
	"sync"

	"github.com/go-kratos/kratos/v2/balancer"
	"github.com/go-kratos/kratos/v2/balancer/node/direct"
	"github.com/go-kratos/kratos/v2/balancer/selector/random"

	"google.golang.org/grpc/attributes"
	gBalancer "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

var (
	_ base.PickerBuilder = &Builder{}
	_ gBalancer.Picker   = &Picker{}

	mu sync.Mutex
)

type node struct {
	balancer.Node
	gBalancer.SubConn
}

func init() {
	SetGlobalBalancer("random", random.New(), &direct.Builder{})
}

// SetGlobalBalancer set grpc balancer with scheme
func SetGlobalBalancer(scheme string, selector balancer.Selector, builder balancer.NodeBuilder) {
	mu.Lock()
	defer mu.Unlock()

	b := base.NewBalancerBuilder(
		scheme,
		&Builder{selector, builder},
		base.Config{HealthCheck: true},
	)
	gBalancer.Register(b)
}

type Builder struct {
	selector    balancer.Selector
	nodeBuilder balancer.NodeBuilder
}

func (b *Builder) Build(info base.PickerBuildInfo) gBalancer.Picker {
	p := &Picker{
		selector: b.selector,
	}
	for conn, info := range info.ReadySCs {
		attr := info.Address.Attributes
		if attr == nil {
			attr = attributes.New()
		}
		p.nodes = append(p.nodes, node{b.nodeBuilder.Build(info.Address.Addr, 100, Attributes(*attr)), conn})
	}
	return p
}

type Picker struct {
	nodes    []balancer.Node
	selector balancer.Selector
}

// Pick pick instances
func (p *Picker) Pick(info gBalancer.PickInfo) (gBalancer.PickResult, error) {
	n, err := p.selector.Select(info.Ctx, p.nodes)
	if err != nil {
		return gBalancer.PickResult{}, err
	}
	done := n.Pick()
	sub := n.(node).SubConn

	return gBalancer.PickResult{SubConn: sub, Done: func(di gBalancer.DoneInfo) {
		done(info.Ctx, balancer.DoneInfo{
			Err:           di.Err,
			BytesSent:     di.BytesSent,
			BytesReceived: di.BytesReceived,
			ReplyHeader:   Trailer(di.Trailer),
		})
	}}, nil
}

type Attributes attributes.Attributes

func (a Attributes) Get(k string) string {
	attr := attributes.Attributes(a)
	v, _ := attr.Value(k).(string)
	return v
}

type Trailer metadata.MD

func (t Trailer) Get(k string) string {
	v := metadata.MD(t).Get(k)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}
