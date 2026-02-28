package selector

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

var _ transport.Transporter = (*Transport)(nil)

type Transport struct {
	kind      transport.Kind
	endpoint  string
	operation string
	headers   *mockHeader
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
	return tr.headers
}

func (tr *Transport) ReplyHeader() transport.Header {
	return nil
}

type mockHeader struct {
	m map[string][]string
}

func (m *mockHeader) Get(key string) string {
	vals := m.m[key]
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (m *mockHeader) Set(key, value string) {
	m.m[key] = []string{value}
}

func (m *mockHeader) Add(key, value string) {
	m.m[key] = append(m.m[key], value)
}

func (m *mockHeader) Keys() []string {
	keys := make([]string, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}
	return keys
}

func (m *mockHeader) Values(key string) []string {
	return m.m[key]
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
			next := func(_ context.Context, req any) (any, error) {
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
			next := func(_ context.Context, req any) (any, error) {
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
			next := func(_ context.Context, req any) (any, error) {
				t.Log(req)
				return "reply", nil
			}
			next = Server(testMiddleware).Match(func(_ context.Context, operation string) bool {
				if strings.HasPrefix(operation, "/go-kratos.dev") || strings.HasSuffix(operation, "world") {
					return true
				}
				return false
			}).Build()(next)
			reply, err := next(test.ctx, test.name)
			if err != nil {
				t.Errorf("expect error is nil, but got %v", err)
			}
			if !reflect.DeepEqual(reply, "reply") {
				t.Errorf("expect reply is reply,but got %v", reply)
			}
		})
	}
}

func TestHeaderFunc(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "/hello.Update/world",
			ctx: transport.NewServerContext(context.Background(), &Transport{
				operation: "/hello.Update/world",
				headers:   &mockHeader{map[string][]string{"X-Test": {"test"}}},
			}),
		},
		{
			name: "/hi.Create/world",
			ctx: transport.NewServerContext(context.Background(), &Transport{
				operation: "/hi.Create/world",
				headers:   &mockHeader{map[string][]string{"X-Test": {"test2"}, "go-kratos": {"kratos"}}},
			}),
		},
		{
			name: "/test.Name/1234",
			ctx: transport.NewServerContext(context.Background(), &Transport{
				operation: "/test.Name/1234",
				headers:   &mockHeader{map[string][]string{"X-Test": {"test3"}}},
			}),
		},
		{
			name: "/go-kratos.dev/kratos",
			ctx: transport.NewServerContext(context.Background(), &Transport{
				operation: "/go-kratos.dev/kratos",
				headers:   &mockHeader{map[string][]string{"X-Test": {"test"}}},
			}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next := func(_ context.Context, req any) (any, error) {
				t.Log(req)
				return "reply", nil
			}
			next = Server(testMiddleware).Match(func(ctx context.Context, _ string) bool {
				tr, ok := transport.FromServerContext(ctx)
				if !ok {
					return false
				}
				if tr.RequestHeader().Get("X-Test") == "test" {
					return true
				}
				if tr.RequestHeader().Get("go-kratos") == "kratos" {
					return true
				}
				return false
			}).Build()(next)
			reply, err := next(test.ctx, test.name)
			if err != nil {
				t.Errorf("expect error is nil, but got %v", err)
			}
			if !reflect.DeepEqual(reply, "reply") {
				t.Errorf("expect reply is reply,but got %v", reply)
			}
		})
	}
}

func testMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req any) (reply any, err error) {
		reply, err = handler(ctx, req)
		return
	}
}

func Test_RegexMatch(t *testing.T) {
	tests := []struct {
		name      string
		regex     []string
		operation string
		want      bool
	}{
		{
			name:      "exact match with digits",
			regex:     []string{`/test/[0-9]+`},
			operation: "/test/1234",
			want:      true,
		},
		{
			name:      "no match",
			regex:     []string{`/test/[0-9]+`},
			operation: "/test/abc",
			want:      false,
		},
		{
			name:      "multiple patterns first matches",
			regex:     []string{`/api/v[0-9]+/.*`, `/test/.*`},
			operation: "/api/v2/users",
			want:      true,
		},
		{
			name:      "multiple patterns second matches",
			regex:     []string{`/api/v[0-9]+/.*`, `/test/.*`},
			operation: "/test/hello",
			want:      true,
		},
		{
			name:      "multiple patterns none match",
			regex:     []string{`/api/v[0-9]+/.*`, `/test/[0-9]+`},
			operation: "/other/path",
			want:      false,
		},
		{
			name:      "invalid regex is skipped",
			regex:     []string{"^\b(?"},
			operation: "something",
			want:      false,
		},
		{
			name:      "invalid regex mixed with valid",
			regex:     []string{"^\b(?", `/test/[0-9]+`},
			operation: "/test/1234",
			want:      true,
		},
		{
			name:      "empty regex list",
			regex:     []string{},
			operation: "/test/1234",
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var middlewareApplied bool
			markMiddleware := func(handler middleware.Handler) middleware.Handler {
				return func(ctx context.Context, req any) (any, error) {
					middlewareApplied = true
					return handler(ctx, req)
				}
			}
			next := func(_ context.Context, _ any) (any, error) {
				return "reply", nil
			}
			ctx := transport.NewServerContext(context.Background(), &Transport{operation: tt.operation})
			handler := Server(markMiddleware).Regex(tt.regex...).Build()(next)
			_, _ = handler(ctx, tt.operation)
			if middlewareApplied != tt.want {
				t.Errorf("middleware applied = %v, want %v", middlewareApplied, tt.want)
			}
		})
	}
}

func Test_InvalidRegexSkipped(t *testing.T) {
	b := Server(testMiddleware).Regex("^\b(?", `/valid/[0-9]+`)
	m := b.Build()
	if m == nil {
		t.Fatal("Build() must not return nil")
	}
	if len(b.compiled) != 1 {
		t.Errorf("expected 1 compiled regex, got %d", len(b.compiled))
	}
}

func Test_matches(t *testing.T) {
	b := Builder{}
	if b.matches(context.Background(), func(_ context.Context) (transport.Transporter, bool) { return nil, false }) {
		t.Error("The matches method must return false.")
	}
}
