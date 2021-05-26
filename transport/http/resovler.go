package http

import (
	"context"
	"net/url"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// Target is resolver target
type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

type resolver struct {
	lock  sync.RWMutex
	nodes []*registry.ServiceInstance

	target  Target
	watcher registry.Watcher
	logger  *log.Helper
}

func newResolver(ctx context.Context, scheme string, discovery registry.Discovery, target Target) (*resolver, error) {
	watcher, err := discovery.Watch(ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}
	r := &resolver{
		target:  target,
		watcher: watcher,
		logger:  log.NewHelper(log.DefaultLogger),
	}
	go func() {
		for {
			services, err := watcher.Next()
			if err != nil {
				r.logger.Errorf("http client watch services got unexpected error:=%v", err)
				return
			}
			var nodes []*registry.ServiceInstance
			for _, in := range services {
				endpoint, err := parseEndpoint(scheme, in.Endpoints)
				if err != nil {
					r.logger.Errorf("Failed to parse discovery endpoint: %v error %v", in.Endpoints, err)
					continue
				}
				if endpoint == "" {
					continue
				}
				nodes = append(nodes, in)
			}
			if len(nodes) != 0 {
				r.lock.Lock()
				r.nodes = nodes
				r.lock.Unlock()
			}
		}
	}()
	return r, nil
}

func (r *resolver) fetch(ctx context.Context) []*registry.ServiceInstance {
	r.lock.RLock()
	nodes := r.nodes
	r.lock.RUnlock()

	return nodes
}

func parseEndpoint(schema string, endpoints []string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}
		if u.Scheme == schema {
			return u.Host, nil
		}
	}
	return "", nil
}
