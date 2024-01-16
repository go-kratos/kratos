package direct

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

const name = "direct"

func init() {
	resolver.Register(NewBuilder())
}

type directBuilder struct{}

// NewBuilder creates a directBuilder which is used to factory direct resolvers.
// example:
//
//	direct://<authority>/127.0.0.1:9000,127.0.0.2:9000
func NewBuilder() resolver.Builder {
	return &directBuilder{}
}

func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}
	return newDirectResolver(cc, addrs), nil
}

func (d *directBuilder) Scheme() string {
	return name
}
