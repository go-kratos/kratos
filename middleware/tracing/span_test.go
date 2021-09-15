package tracing

import (
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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
