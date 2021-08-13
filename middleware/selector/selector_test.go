package selector

import (
	"context"
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

func TestMatch(t *testing.T) {

	tests := []struct {
		name string
		ctx  context.Context
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "/hello/world",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/hello/world"}),
			want: true,
		},
		{
			name: "/hi/world",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/hi/world"}),
			want: false,
		},
		{
			name: "/test/1234",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/test/1234"}),
			want: true,
		},
		{
			name: "/example/kratos",
			ctx:  transport.NewServerContext(context.Background(), &Transport{operation: "/example/kratos"}),
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				// t.Log(req)
				return "reply", nil
			}
			var marker markerMiddleware
			next = Server(marker.middleware).
				Prefix("/hello/").
				Regex(`/test/[0-9]+`).
				Path("/example/kratos").
				Build()(next)
			next(test.ctx, test.name)
			if got := marker.marked; got != test.want {
				t.Errorf("Match() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestMatchClient(t *testing.T) {

	tests := []struct {
		name string
		ctx  context.Context
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "/hello/world",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/hello/world"}),
			want: true,
		},
		{
			name: "/hi/world",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/hi/world"}),
			want: false,
		},
		{
			name: "/test/1234",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/test/1234"}),
			want: true,
		},
		{
			name: "/example/kratos",
			ctx:  transport.NewClientContext(context.Background(), &Transport{operation: "/example/kratos"}),
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				// t.Log(req)
				return "reply", nil
			}
			var marker markerMiddleware
			next = Client(marker.middleware).
				Prefix("/hello/").
				Regex(`/test/[0-9]+`).
				Path("/example/kratos").
				Build()(next)
			next(test.ctx, test.name)
			if got := marker.marked; got != test.want {
				t.Errorf("Match() = %v, want %v", got, test.want)
			}
		})
	}
}

type markerMiddleware struct {
	marked bool
}

func (m *markerMiddleware) middleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		m.marked = true
		// fmt.Println("before")
		reply, err = handler(ctx, req)
		// fmt.Println("after")
		return
	}
}

func TestBuilder_opMatcher_Match(t *testing.T) {

	builder := Server().
		Prefix("/hello/").
		Regex(`/test/[0-9]+`).
		Path("/example/kratos")

	opMatcher := builder.rootOpMatcherBuilder.build()

	tests := []struct {
		name      string
		operation string
		want      bool
	}{
		{operation: "/hello/world", want: true},
		{operation: "/hi/world", want: false},
		{operation: "/test/1234", want: true},
		{operation: "/example/kratos", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := opMatcher.Match(tt.operation); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
