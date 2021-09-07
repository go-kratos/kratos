package selector

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/stretchr/testify/assert"
)

var _ transport.Transporter = &Transport{}

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
}

func (tr *Transport) Kind() transport.Kind {
	return tr.kind
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) Operation() string {
	return tr.operation
}

func (tr *Transport) RequestHeader() transport.Header {
	return nil
}

func (tr *Transport) ReplyHeader() transport.Header {
	return nil
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		// TODO: Add test cases.
		{
			name: "/hello/world",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/hello/world"}),
		},
		{
			name: "/hi/world",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/hi/world"}),
		},
		{
			name: "/test/1234",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/test/1234"}),
		},
		{
			name: "/example/kratos",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/example/kratos"}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				return "reply", nil
			}
			next = Server(testMiddleware).Prefix("/hello/").Regex(`/test/[0-9]+`).
				Path("/example/kratos").Build()(next)
			_, _ = next(test.ctx, test.name)
		})
	}
}

func TestMatchClient(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		// TODO: Add test cases.
		{
			name: "/hello/world",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/hello/world"}),
		},
		{
			name: "/hi/world",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/hi/world"}),
		},
		{
			name: "/test/1234",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/test/1234"}),
		},
		{
			name: "/example/kratos",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/example/kratos"}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				return "reply", nil
			}
			next = Client(testMiddleware).Prefix("/hello/").Regex(`/test/[0-9]+`).
				Path("/example/kratos").Build()(next)
			_, _ = next(test.ctx, test.name)
		})
	}
}

func TestFunc(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "/hello.Update/world",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/hello.Update/world"}),
		},
		{
			name: "/hi.Create/world",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/hi.Create/world"}),
		},
		{
			name: "/test.Name/1234",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/test.Name/1234"}),
		},
		{
			name: "/go-kratos.dev/kratos",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/go-kratos.dev/kratos"}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				return "reply", nil
			}
			next = Server(testMiddleware).Match(func(operation string) bool {
				if strings.HasPrefix(operation, "/go-kratos.dev") || strings.HasSuffix(operation, "world") {
					return true
				}
				return false
			}).Build()(next)
			reply, err := next(test.ctx, test.name)
			assert.Equal(t, reply, "reply")
			assert.Nil(t, err)
		})
	}
}

func testMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("before")
		reply, err = handler(ctx, req)
		fmt.Println("after")
		return
	}
}
