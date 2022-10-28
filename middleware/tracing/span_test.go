package tracing

import (
	"context"
	"net"
	"net/http"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/peer"

	"github.com/go-kratos/kratos/v2/internal/testdata/binding"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
)

func Test_parseFullMethod(t *testing.T) {
	tests := []struct {
		name       string
		fullMethod string
		want       string
		wantAttr   []attribute.KeyValue
	}{
		{
			name:       "/foo.bar/hello",
			fullMethod: "/foo.bar/hello",
			want:       "foo.bar/hello",
			wantAttr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("foo.bar"),
				semconv.RPCMethodKey.String("hello"),
			},
		},
		{
			name:       "/foo.bar/hello/world",
			fullMethod: "/foo.bar/hello/world",
			want:       "foo.bar/hello/world",
			wantAttr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("foo.bar"),
				semconv.RPCMethodKey.String("hello/world"),
			},
		},
		{
			name:       "/hello",
			fullMethod: "/hello",
			want:       "hello",
			wantAttr:   []attribute.KeyValue{attribute.Key("rpc.operation").String("/hello")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseFullMethod(tt.fullMethod)
			if got != tt.want {
				t.Errorf("parseFullMethod() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantAttr) {
				t.Errorf("parseFullMethod() got1 = %v, want %v", got1, tt.wantAttr)
			}
		})
	}
}

func Test_peerAttr(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want []attribute.KeyValue
	}{
		{
			name: "nil addr",
			addr: ":8080",
			want: []attribute.KeyValue{
				semconv.NetPeerIPKey.String("127.0.0.1"),
				semconv.NetPeerPortKey.String("8080"),
			},
		},
		{
			name: "normal addr without port",
			addr: "192.168.0.1",
			want: []attribute.KeyValue(nil),
		},
		{
			name: "normal addr with port",
			addr: "192.168.0.1:8080",
			want: []attribute.KeyValue{
				semconv.NetPeerIPKey.String("192.168.0.1"),
				semconv.NetPeerPortKey.String("8080"),
			},
		},
		{
			name: "dns addr",
			addr: "foo:8080",
			want: []attribute.KeyValue{
				semconv.NetPeerIPKey.String("foo"),
				semconv.NetPeerPortKey.String("8080"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := peerAttr(tt.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("peerAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTarget(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		wantAddress string
		wantErr     bool
	}{
		{
			name:        "http",
			endpoint:    "http://foo.bar:8080",
			wantAddress: "http://foo.bar:8080",
			wantErr:     false,
		},
		{
			name:        "http",
			endpoint:    "http://127.0.0.1:8080",
			wantAddress: "http://127.0.0.1:8080",
			wantErr:     false,
		},
		{
			name:        "without protocol",
			endpoint:    "foo.bar:8080",
			wantAddress: "foo.bar:8080",
			wantErr:     false,
		},
		{
			name:        "grpc",
			endpoint:    "grpc://foo.bar",
			wantAddress: "grpc://foo.bar",
			wantErr:     false,
		},
		{
			name:        "with path",
			endpoint:    "/foo",
			wantAddress: "foo",
			wantErr:     false,
		},
		{
			name:        "with path",
			endpoint:    "http://127.0.0.1/hello",
			wantAddress: "hello",
			wantErr:     false,
		},
		{
			name:        "empty",
			endpoint:    "%%",
			wantAddress: "",
			wantErr:     true,
		},
		{
			name:        "invalid path",
			endpoint:    "//%2F/#%2Fanother",
			wantAddress: "",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, err := parseTarget(tt.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("parseTarget() = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}

func Test_setServerSpan(t *testing.T) {
	ctx := context.Background()
	_, span := trace.NewNoopTracerProvider().Tracer("Tracer").Start(ctx, "Spanname")

	// Handle without Transport context
	setServerSpan(ctx, span, nil)

	// Handle with proto message
	m := &binding.HelloRequest{}
	setServerSpan(ctx, span, m)

	// Handle with metadata context
	ctx = metadata.NewServerContext(ctx, metadata.New())
	setServerSpan(ctx, span, m)

	// Handle with KindHTTP transport context
	mt := &mockTransport{
		kind: transport.KindHTTP,
	}
	mt.request, _ = http.NewRequest(http.MethodGet, "/endpoint", nil)
	ctx = transport.NewServerContext(ctx, mt)
	setServerSpan(ctx, span, m)

	// Handle with KindGRPC transport context
	mt.kind = transport.KindGRPC
	ctx = transport.NewServerContext(ctx, mt)
	ip, _ := net.ResolveIPAddr("ip", "1.1.1.1")
	ctx = peer.NewContext(ctx, &peer.Peer{
		Addr: ip,
	})
	setServerSpan(ctx, span, m)
}

func Test_setClientSpan(t *testing.T) {
	ctx := context.Background()
	_, span := trace.NewNoopTracerProvider().Tracer("Tracer").Start(ctx, "Spanname")

	// Handle without Transport context
	setClientSpan(ctx, span, nil)

	// Handle with proto message
	m := &binding.HelloRequest{}
	setClientSpan(ctx, span, m)

	// Handle with metadata context
	ctx = metadata.NewClientContext(ctx, metadata.New())
	setClientSpan(ctx, span, m)

	// Handle with KindHTTP transport context
	mt := &mockTransport{
		kind: transport.KindHTTP,
	}
	mt.request, _ = http.NewRequest(http.MethodGet, "/endpoint", nil)
	mt.request.Host = "MyServer"
	ctx = transport.NewClientContext(ctx, mt)
	setClientSpan(ctx, span, m)

	// Handle with KindGRPC transport context
	mt.kind = transport.KindGRPC
	ctx = transport.NewClientContext(ctx, mt)
	ip, _ := net.ResolveIPAddr("ip", "1.1.1.1")
	ctx = peer.NewContext(ctx, &peer.Peer{
		Addr: ip,
	})
	setClientSpan(ctx, span, m)

	// Handle without Host request
	ctx = transport.NewClientContext(ctx, mt)
	setClientSpan(ctx, span, m)
}
