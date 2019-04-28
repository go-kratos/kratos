package direct

import (
	"context"
	"fmt"
	"strings"

	"go-common/library/conf/env"
	"go-common/library/naming"
	"go-common/library/net/rpc/warden/resolver"
)

const (
	// Name is the name of direct resolver
	Name = "direct"
)

var _ naming.Resolver = &Direct{}

// New return Direct
func New() *Direct {
	return &Direct{}
}

// Build build direct.
func Build(id string) *Direct {
	return &Direct{id: id}
}

// Direct is a resolver for conneting endpoints directly.
// example format: direct://default/192.168.1.1:8080,192.168.1.2:8081
type Direct struct {
	id string
}

// Build direct build.
func (d *Direct) Build(id string) naming.Resolver {
	return &Direct{id: id}
}

// Scheme return the Scheme of Direct
func (d *Direct) Scheme() string {
	return Name
}

// Watch a tree
func (d *Direct) Watch() <-chan struct{} {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return ch
}

//Unwatch a tree
func (d *Direct) Unwatch(id string) {
}

//Fetch fetch isntances
func (d *Direct) Fetch(ctx context.Context) (insMap map[string][]*naming.Instance, found bool) {
	var ins []*naming.Instance

	addrs := strings.Split(d.id, ",")
	for _, addr := range addrs {
		ins = append(ins, &naming.Instance{
			Addrs: []string{fmt.Sprintf("%s://%s", resolver.Scheme, addr)},
		})
	}
	if len(ins) > 0 {
		found = true
	}
	insMap = map[string][]*naming.Instance{env.Zone: ins}
	return
}

//Close close Direct
func (d *Direct) Close() error {
	return nil
}
