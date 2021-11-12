package filter

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, "127.0.0.2:9090", nodes[0].Address())
}
