package grpc

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestTrailer(t *testing.T) {
	trailer := Trailer(metadata.New(map[string]string{"a": "b"}))
	assert.Equal(t, "b", trailer.Get("a"))
	assert.Equal(t, "", trailer.Get("3"))
}

func TestBalancerName(t *testing.T) {
	o := &clientOptions{}

	WithBalancerName("p2c")(o)
	assert.Equal(t, "p2c", o.balancerName)
}

func TestFilters(t *testing.T) {
	o := &clientOptions{}

	WithFilter(func(_ context.Context, nodes []selector.Node) []selector.Node {
		return nodes
	})(o)
	assert.Equal(t, 1, len(o.filters))
}
