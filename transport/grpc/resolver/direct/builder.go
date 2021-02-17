package direct

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(NewBuilder())
}

type directBuilder struct{}

// NewBuilder creates a directBuilder which is used to factory direct resolvers.
// example:
//   direct://<authority>/127.0.0.1:9000,127.0.0.2:9000
func NewBuilder() resolver.Builder {
	return &directBuilder{}
}

func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var addrs []resolver.Address
	for _, addr := range strings.Split(target.Endpoint, ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}
	cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
	return newDirectResolver(), nil
}

func (d *directBuilder) Scheme() string {
	return "direct"
}
