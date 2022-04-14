package http

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/internal/endpoint"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
)

// Target is resolver target
type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

func parseTarget(endpoint string, insecure bool) (*Target, error) {
	if !strings.Contains(endpoint, "://") {
		if insecure {
			endpoint = "http://" + endpoint
		} else {
			endpoint = "https://" + endpoint
		}
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	target := &Target{Scheme: u.Scheme, Authority: u.Host}
	if len(u.Path) > 1 {
		target.Endpoint = u.Path[1:]
	}
	return target, nil
}

type resolver struct {
	rebalancer selector.Rebalancer

	target  *Target
	watcher registry.Watcher
	logger  *log.Helper

	insecure bool
}

func newResolver(ctx context.Context, discovery registry.Discovery, target *Target, rebalancer selector.Rebalancer, block, insecure bool) (*resolver, error) {
	watcher, err := discovery.Watch(ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}
	r := &resolver{
		target:     target,
		watcher:    watcher,
		logger:     log.NewHelper(log.GetLogger()),
		rebalancer: rebalancer,
		insecure:   insecure,
	}
	if block {
		done := make(chan error, 1)
		go func() {
			for {
				services, err := watcher.Next()
				if err != nil {
					done <- err
					return
				}
				if r.update(services) {
					done <- nil
					return
				}
			}
		}()
		select {
		case err := <-done:
			if err != nil {
				err := watcher.Stop()
				if err != nil {
					r.logger.Errorf("failed to http client watch stop: %v", target)
				}
				return nil, err
			}
		case <-ctx.Done():
			r.logger.Errorf("http client watch service %v reaching context deadline!", target)
			err := watcher.Stop()
			if err != nil {
				r.logger.Errorf("failed to http client watch stop: %v", target)
			}
			return nil, ctx.Err()
		}
	}
	go func() {
		for {
			services, err := watcher.Next()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				r.logger.Errorf("http client watch service %v got unexpected error:=%v", target, err)
				time.Sleep(time.Second)
				continue
			}
			r.update(services)
		}
	}()
	return r, nil
}

func (r *resolver) update(services []*registry.ServiceInstance) bool {
	nodes := make([]selector.Node, 0)
	for _, ins := range services {
		ept, err := endpoint.ParseEndpoint(ins.Endpoints, "http", !r.insecure)
		if err != nil {
			r.logger.Errorf("Failed to parse (%v) discovery endpoint: %v error %v", r.target, ins.Endpoints, err)
			continue
		}
		if ept == "" {
			continue
		}
		nodes = append(nodes, selector.NewNode("http", ept, ins))
	}
	if len(nodes) == 0 {
		r.logger.Warnf("[http resolver]Zero endpoint found,refused to write,set: %s ins: %v", r.target.Endpoint, nodes)
		return false
	}
	r.rebalancer.Apply(nodes)
	return true
}

func (r *resolver) Close() error {
	return r.watcher.Stop()
}
