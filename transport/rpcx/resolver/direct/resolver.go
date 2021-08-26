package direct

import "google.golang.org/grpc/resolver"

type directResolver struct{}

func newDirectResolver() resolver.Resolver {
	return &directResolver{}
}

func (r *directResolver) Close() {
}

func (r *directResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
