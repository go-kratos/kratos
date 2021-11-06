package filter

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
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

	f(context.Background(), &nodes)
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, "127.0.0.2:9090", nodes[0].Address())
}

func TestBaseFilterWithRandom(t *testing.T) {
	for i := 0; i < 100; i++ {
		testBaseFilter(t, 1000, rand.Intn(1000))
	}

	testBaseFilter(t, 0, rand.Intn(1000))
	testBaseFilter(t, 1, 1000)
	testBaseFilter(t, 2, 1000)
	testBaseFilter(t, 3, 1000)
	testBaseFilter(t, 1, 0)
	testBaseFilter(t, 2, 0)
	testBaseFilter(t, 3, 0)
}

func testBaseFilter(t *testing.T, length int, reservedRatio int) {
	var raw []selector.Node
	var targets map[string]selector.Node = make(map[string]selector.Node)
	for i := 0; i < length; i++ {
		addr := strconv.FormatInt(int64(i), 10)
		raw = append(raw, selector.NewNode(
			addr,
			&registry.ServiceInstance{
				ID:        addr,
				Name:      "helloworld",
				Endpoints: []string{addr},
			}))
		if reservedRatio > rand.Intn(length) {
			targets[addr] = raw[i]
		}
	}

	keepFactory := func(ctx context.Context) Keep {
		return func(node selector.Node) bool {
			if _, ok := targets[node.Address()]; ok {
				return true
			}
			return false
		}
	}
	f := BaseFilter(keepFactory)
	f(context.Background(), &raw)
	assert.Equal(t, len(targets), len(raw))
	for _, n := range raw {
		_, ok := targets[n.Address()]
		assert.True(t, ok)
	}
	fmt.Println(len(targets))
}
