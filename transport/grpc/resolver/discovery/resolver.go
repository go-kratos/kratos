package discovery

import (
	"context"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

type discoveryResolver struct {
	w   registry.Watcher
	cc  resolver.ClientConn
	log *log.Helper

	ctx    context.Context
	cancel context.CancelFunc
}

func (r *discoveryResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}

		ins, err := r.w.Next()
		if err != nil {
			r.log.Errorf("Failed to watch discovery endpoint: %v", err)
			time.Sleep(time.Second)
			continue
		}
		r.update(ins)
	}
}

func (r *discoveryResolver) update(ins []*registry.ServiceInstance) {
	var addrs []resolver.Address
	for _, in := range ins {
		endpoint, err := parseEndpoint(in.Endpoints)
		if err != nil {
			r.log.Errorf("Failed to parse discovery endpoint: %v", err)
			continue
		}
		if endpoint == "" {
			continue
		}
		addr := resolver.Address{
			ServerName: in.Name,
			Attributes: parseAttributes(in.Metadata),
			Addr:       endpoint,
		}
		addrs = append(addrs, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *discoveryResolver) Close() {
	r.cancel()
	r.w.Stop()
}

func (r *discoveryResolver) ResolveNow(options resolver.ResolveNowOptions) {}

func parseEndpoint(endpoints []string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}
		if u.Scheme == "grpc" {
			return u.Host, nil
		}
	}
	return "", nil
}

func parseAttributes(md map[string]string) *attributes.Attributes {
	pairs := make([]interface{}, 0, len(md))
	for k, v := range md {
		pairs = append(pairs, k, v)
	}
	return attributes.New(pairs...)
}
