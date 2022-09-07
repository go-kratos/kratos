package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/selector"
	"google.golang.org/grpc/metadata"
)

func TestTrailer(t *testing.T) {
	trailer := Trailer(metadata.New(map[string]string{"a": "b"}))
	if !reflect.DeepEqual("b", trailer.Get("a")) {
		t.Errorf("expect %v, got %v", "b", trailer.Get("a"))
	}
	if !reflect.DeepEqual("", trailer.Get("notfound")) {
		t.Errorf("expect %v, got %v", "", trailer.Get("notfound"))
	}
}

func TestFilters(t *testing.T) {
	o := &clientOptions{}

	WithNodeFilter(func(_ context.Context, nodes []selector.Node) []selector.Node {
		return nodes
	})(o)
	if !reflect.DeepEqual(1, len(o.filters)) {
		t.Errorf("expect %v, got %v", 1, len(o.filters))
	}
}
