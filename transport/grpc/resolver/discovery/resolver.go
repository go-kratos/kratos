package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"

	"github.com/go-kratos/kratos/v3/internal/endpoint"
	"github.com/go-kratos/kratos/v3/internal/subset"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/registry"
)

type discoveryResolver struct {
	w  registry.Watcher
	cc resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	insecure    bool
	selectorKey string
	subsetSize  int
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
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Error("[resolver] failed to watch discovery endpoint", "error", err)
			time.Sleep(time.Second)
			continue
		}
		r.update(ins)
	}
}

func (r *discoveryResolver) update(ins []*registry.ServiceInstance) {
	var (
		endpoints = make(map[string]struct{})
		filtered  = make([]*registry.ServiceInstance, 0, len(ins))
	)
	for _, in := range ins {
		ept, err := endpoint.ParseEndpoint(in.Endpoints, endpoint.Scheme("grpc", !r.insecure))
		if err != nil {
			log.Error("[resolver] failed to parse discovery endpoint", "error", err)
			continue
		}
		if ept == "" {
			continue
		}
		// filter redundant endpoints
		if _, ok := endpoints[ept]; ok {
			continue
		}
		endpoints[ept] = struct{}{}
		filtered = append(filtered, in)
	}
	if r.subsetSize != 0 {
		filtered = subset.Subset(r.selectorKey, filtered, r.subsetSize)
	}

	addrs := make([]resolver.Address, 0, len(filtered))
	for _, in := range filtered {
		ept, _ := endpoint.ParseEndpoint(in.Endpoints, endpoint.Scheme("grpc", !r.insecure))
		addr := resolver.Address{
			ServerName: in.Name,
			Attributes: parseAttributes(in.Metadata).WithValue("rawServiceInstance", in),
			Addr:       ept,
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		log.Warn("[resolver] zero endpoint found, refused to write", "instances", ins)
		return
	}
	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		log.Error("[resolver] failed to update state", "error", err)
	}

	b, _ := json.Marshal(filtered)
	log.Info("[resolver] update instances", "instances", string(b))
}

func (r *discoveryResolver) Close() {
	r.cancel()
	err := r.w.Stop()
	if err != nil {
		log.Error("[resolver] failed to stop watcher", "error", err)
	}
}

func (r *discoveryResolver) ResolveNow(_ resolver.ResolveNowOptions) {}

func parseAttributes(md map[string]string) (a *attributes.Attributes) {
	for k, v := range md {
		a = a.WithValue(k, v)
	}
	return a
}
