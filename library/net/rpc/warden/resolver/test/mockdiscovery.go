package resolver

import (
	"context"
	"go-common/library/conf/env"
	"go-common/library/naming"
)

type mockDiscoveryBuilder struct {
	instances map[string]*naming.Instance
	watchch   map[string][]*mockDiscoveryResolver
}

func (mb *mockDiscoveryBuilder) Build(id string) naming.Resolver {
	mr := &mockDiscoveryResolver{
		d:       mb,
		watchch: make(chan struct{}, 1),
	}
	mb.watchch[id] = append(mb.watchch[id], mr)
	mr.watchch <- struct{}{}
	return mr
}
func (mb *mockDiscoveryBuilder) Scheme() string {
	return "mockdiscovery"
}

type mockDiscoveryResolver struct {
	//instances map[string]*naming.Instance
	d       *mockDiscoveryBuilder
	watchch chan struct{}
}

var _ naming.Resolver = &mockDiscoveryResolver{}

func (md *mockDiscoveryResolver) Fetch(ctx context.Context) (map[string][]*naming.Instance, bool) {
	zones := make(map[string][]*naming.Instance)
	for _, v := range md.d.instances {
		zones[v.Zone] = append(zones[v.Zone], v)
	}
	return zones, len(zones) > 0
}

func (md *mockDiscoveryResolver) Watch() <-chan struct{} {
	return md.watchch
}

func (md *mockDiscoveryResolver) Close() error {
	close(md.watchch)
	return nil
}

func (md *mockDiscoveryResolver) Scheme() string {
	return "mockdiscovery"
}

func (mb *mockDiscoveryBuilder) registry(appID string, hostname, rpc string, metadata map[string]string) {
	ins := &naming.Instance{
		AppID:    appID,
		Env:      "hello=world",
		Hostname: hostname,
		Addrs:    []string{"grpc://" + rpc},
		Version:  "1.1",
		Zone:     env.Zone,
		Metadata: metadata,
	}
	mb.instances[hostname] = ins
	if ch, ok := mb.watchch[appID]; ok {
		var bullet struct{}
		for _, c := range ch {
			c.watchch <- bullet
		}
	}
}

func (mb *mockDiscoveryBuilder) cancel(hostname string) {
	ins, ok := mb.instances[hostname]
	if !ok {
		return
	}
	delete(mb.instances, hostname)
	if ch, ok := mb.watchch[ins.AppID]; ok {
		var bullet struct{}
		for _, c := range ch {
			c.watchch <- bullet
		}
	}
}
