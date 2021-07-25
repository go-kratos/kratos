package selector

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"testing"
)

var (
	_ transport.Transporter = &Transport{}
)

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

func TestMatchFull(t *testing.T) {
	type args struct {
		route string
		ms    []middleware.Middleware
	}
	tests := []struct {
		name string
		args args
		ctx  context.Context
	}{
		// TODO: Add test cases.
		{
			name: "/hello/world",
			args: args{
				route: "/hello/world",
				ms:    []middleware.Middleware{testMiddleware},
			},
			ctx: transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/hello/world"}),
		},
		{
			name: "/hello/world/test",
			args: args{
				route: "/hello/world",
				ms:    []middleware.Middleware{testMiddleware},
			},
			ctx: transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/hello"}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				return "reply", nil
			}
			next = MatchFull(test.args.route, test.args.ms...)(next)
			next(test.ctx, test.name)
		})
	}
}

func TestMatchPrefix(t *testing.T) {
	type args struct {
		prefix string
		ms     []middleware.Middleware
	}
	tests := []struct {
		name string
		args args
		ctx  context.Context
	}{
		// TODO: Add test cases.
		{
			name: "/hello/world",
			args: args{
				prefix: "/hello/",
				ms:     []middleware.Middleware{testMiddleware},
			},
			ctx: transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/hello/world"}),
		},
		{
			name: "/hi/world",
			args: args{
				prefix: "/hello/",
				ms:     []middleware.Middleware{testMiddleware},
			},
			ctx: transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/hi/world"}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				return "reply", nil
			}
			next = MatchPrefix(test.args.prefix, test.args.ms...)(next)
			next(test.ctx, test.name)
		})
	}
}

func TestMatchRegex(t *testing.T) {
	type args struct {
		pattern string
		ms      []middleware.Middleware
	}
	tests := []struct {
		name string
		args args
		ctx  context.Context
	}{
		// TODO: Add test cases.
		{
			name: "/hello/1234",
			args: args{
				pattern: `/hello/[0-9]+`,
				ms:      []middleware.Middleware{testMiddleware},
			},
			ctx: transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/hello/1234"}),
		},
		{
			name: "/hello/test",
			args: args{
				pattern: `/hello/[0-9]+`,
				ms:      []middleware.Middleware{testMiddleware},
			},
			ctx: transport.NewServerContext(context.Background(), &Transport{kind: transport.KindHTTP, endpoint: "endpoint", operation: "/hello/test"}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				t.Log(req)
				return "reply", nil
			}
			next = MatchRegex(test.args.pattern, test.args.ms...)(next)
			next(test.ctx, test.name)
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
