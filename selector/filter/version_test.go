package filter

import (
	"context"
	"reflect"
	"testing"

	"github.com/SeeMusic/kratos/v2/registry"
	"github.com/SeeMusic/kratos/v2/selector"
)

func TestVersion(t *testing.T) {
	f := Version("v2.0.0")
	var nodes []selector.Node
	nodes = append(nodes, selector.NewNode(
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
		}))

	nodes = append(nodes, selector.NewNode(
		"127.0.0.2:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.2:9090",
			Name:      "helloworld",
			Version:   "v2.0.0",
			Endpoints: []string{"http://127.0.0.2:9090"},
		}))

	nodes = f(context.Background(), nodes)
	if !reflect.DeepEqual(len(nodes), 1) {
		t.Errorf("expect %v, got %v", 1, len(nodes))
	}
	if !reflect.DeepEqual(nodes[0].Address(), "127.0.0.2:9090") {
		t.Errorf("expect %v, got %v", nodes[0].Address(), "127.0.0.2:9090")
	}
}
