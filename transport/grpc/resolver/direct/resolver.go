package direct

import (
	"google.golang.org/grpc/resolver"

	"github.com/go-kratos/kratos/v2/log"
)

type directResolver struct {
	addresses []resolver.Address
	cc        resolver.ClientConn
}

func newDirectResolver(cc resolver.ClientConn, addresses []resolver.Address) resolver.Resolver {
	d := &directResolver{cc: cc, addresses: addresses}
	d.update()
	return d
}

func (r *directResolver) Close() {
}

func (r *directResolver) ResolveNow(_ resolver.ResolveNowOptions) {
	r.update()
}

func (r *directResolver) update() {
	err := r.cc.UpdateState(resolver.State{
		Addresses: r.addresses,
	})
	if err != nil {
		log.Errorf("[resolver] failed to update state: %s", err)
	}
}
