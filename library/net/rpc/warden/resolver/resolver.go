package resolver

import (
	"context"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"

	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/naming"
	wmeta "go-common/library/net/rpc/warden/metadata"

	"github.com/dgryski/go-farm"
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
)

const (
	// Scheme is the scheme of discovery address
	Scheme = "grpc"
)

var (
	_  resolver.Resolver = &Resolver{}
	_  resolver.Builder  = &Builder{}
	mu sync.Mutex
)

// Register register resolver builder if nil.
func Register(b naming.Builder) {
	mu.Lock()
	defer mu.Unlock()
	if resolver.Get(b.Scheme()) == nil {
		resolver.Register(&Builder{b})
	}
}

// Set override any registered builder
func Set(b naming.Builder) {
	mu.Lock()
	defer mu.Unlock()
	resolver.Register(&Builder{b})
}

// Builder is also a resolver builder.
// It's build() function always returns itself.
type Builder struct {
	naming.Builder
}

// Build returns itself for Resolver, because it's both a builder and a resolver.
func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var zone = env.Zone
	ss := int64(50)
	clusters := map[string]struct{}{}
	str := strings.SplitN(target.Endpoint, "?", 2)
	if len(str) == 0 {
		return nil, errors.Errorf("warden resolver: parse target.Endpoint(%s) failed!err:=endpoint is empty", target.Endpoint)
	} else if len(str) == 2 {
		m, err := url.ParseQuery(str[1])
		if err == nil {
			for _, c := range m[naming.MetaCluster] {
				clusters[c] = struct{}{}
			}
			zones := m[naming.MetaZone]
			if len(zones) > 0 {
				zone = zones[0]
			}
			if sub, ok := m["subset"]; ok {
				if t, err := strconv.ParseInt(sub[0], 10, 64); err == nil {
					ss = t
				}

			}
		}
	}
	r := &Resolver{
		nr:         b.Builder.Build(str[0]),
		cc:         cc,
		quit:       make(chan struct{}, 1),
		clusters:   clusters,
		zone:       zone,
		subsetSize: ss,
	}
	go r.updateproc()
	return r, nil
}

// Resolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type Resolver struct {
	nr   naming.Resolver
	cc   resolver.ClientConn
	quit chan struct{}

	clusters   map[string]struct{}
	zone       string
	subsetSize int64
}

// Close is a noop for Resolver.
func (r *Resolver) Close() {
	select {
	case r.quit <- struct{}{}:
		r.nr.Close()
	default:
	}
}

// ResolveNow is a noop for Resolver.
func (r *Resolver) ResolveNow(o resolver.ResolveNowOption) {
}

func (r *Resolver) updateproc() {
	event := r.nr.Watch()
	for {
		select {
		case <-r.quit:
			return
		case _, ok := <-event:
			if !ok {
				return
			}
		}
		if insMap, ok := r.nr.Fetch(context.Background()); ok {
			instances, ok := insMap[r.zone]
			if !ok {
				for _, value := range insMap {
					instances = append(instances, value...)
				}
			}
			if r.subsetSize > 0 && len(instances) > 0 {
				instances = r.subset(instances, env.Hostname, r.subsetSize)
			}
			r.newAddress(instances)
		}
	}
}

func (r *Resolver) subset(backends []*naming.Instance, clientID string, size int64) []*naming.Instance {
	if len(backends) <= int(size) {
		return backends
	}
	sort.Slice(backends, func(i, j int) bool {
		return backends[i].Hostname < backends[j].Hostname
	})
	count := int64(len(backends)) / size

	id := farm.Fingerprint64([]byte(clientID))
	round := int64(id / uint64(count))

	s := rand.NewSource(round)
	ra := rand.New(s)
	ra.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})
	start := (id % uint64(count)) * uint64(size)
	return backends[int(start) : int(start)+int(size)]
}

func (r *Resolver) newAddress(instances []*naming.Instance) {
	if len(instances) <= 0 {
		return
	}
	addrs := make([]resolver.Address, 0, len(instances))
	for _, ins := range instances {
		if len(r.clusters) > 0 {
			if _, ok := r.clusters[ins.Metadata[naming.MetaCluster]]; !ok {
				continue
			}
		}

		var weight int64
		if weight, _ = strconv.ParseInt(ins.Metadata[naming.MetaWeight], 10, 64); weight <= 0 {
			weight = 10
		}
		var rpc string
		for _, a := range ins.Addrs {
			u, err := url.Parse(a)
			if err == nil && u.Scheme == Scheme {
				rpc = u.Host
			}
		}
		if rpc == "" {
			log.Warn("warden/resolver: invalid rpc address(%s,%s,%v) found!", ins.AppID, ins.Hostname, ins.Addrs)
			continue
		}
		addr := resolver.Address{
			Addr:       rpc,
			Type:       resolver.Backend,
			ServerName: ins.AppID,
			Metadata:   wmeta.MD{Weight: weight, Color: ins.Metadata[naming.MetaColor]},
		}
		addrs = append(addrs, addr)
	}
	r.cc.NewAddress(addrs)
}
