package metadata

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v3/metadata"
	"github.com/go-kratos/kratos/v3/transport"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string { return http.Header(hc).Get(key) }

func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }

func (hc headerCarrier) Add(key string, value string) { http.Header(hc).Add(key, value) }

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice value associated with the passed key.
func (hc headerCarrier) Values(key string) []string {
	return http.Header(hc).Values(key)
}

type testTransport struct{ header headerCarrier }

func (tr *testTransport) Kind() transport.Kind            { return transport.KindHTTP }
func (tr *testTransport) Endpoint() string                { return "" }
func (tr *testTransport) Operation() string               { return "" }
func (tr *testTransport) RequestHeader() transport.Header { return tr.header }
func (tr *testTransport) ReplyHeader() transport.Header   { return tr.header }

var (
	globalKey   = "x-md-global-key"
	globalValue = "global-value"
	localKey    = "x-md-local-key"
	localValue  = "local-value"
	customKey   = "x-md-local-custom"
	customValue = "custom-value"
	constKey    = "x-md-local-const"
	constValue  = "x-md-local-const"
)

func TestSever(t *testing.T) {
	hs := func(ctx context.Context, in any) (any, error) {
		md, ok := metadata.FromServerContext(ctx)
		if !ok {
			return nil, errors.New("no md")
		}
		if md.Get(constKey) != constValue {
			return nil, errors.New("const not equal")
		}
		if md.Get(globalKey) != globalValue {
			return nil, errors.New("global not equal")
		}
		if md.Get(localKey) != localValue {
			return nil, errors.New("local not equal")
		}
		return in, nil
	}
	hc := headerCarrier{}
	hc.Set(globalKey, globalValue)
	hc.Set(localKey, localValue)
	ctx := transport.NewServerContext(context.Background(), &testTransport{hc})
	// const md
	constMD := metadata.New()
	constMD.Set(constKey, constValue)
	reply, err := Server(WithConstants(constMD))(hs)(ctx, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if reply.(string) != "foo" {
		t.Fatalf("want foo got %v", reply)
	}
}

func TestClient(t *testing.T) {
	hs := func(ctx context.Context, in any) (any, error) {
		tr, ok := transport.FromClientContext(ctx)
		if !ok {
			return nil, errors.New("no md")
		}
		if tr.RequestHeader().Get(constKey) != constValue {
			return nil, errors.New("const not equal")
		}
		if tr.RequestHeader().Get(customKey) != customValue {
			return nil, errors.New("custom not equal")
		}
		if tr.RequestHeader().Get(globalKey) != globalValue {
			return nil, errors.New("global not equal")
		}
		if tr.RequestHeader().Get(localKey) != "" {
			return nil, errors.New("local must empty")
		}
		return in, nil
	}
	// server md
	serverMD := metadata.New()
	serverMD.Set(globalKey, globalValue)
	serverMD.Set(localKey, localValue)
	ctx := metadata.NewServerContext(context.Background(), serverMD)
	// client md
	clientMD := metadata.New()
	clientMD.Set(customKey, customValue)
	ctx = metadata.NewClientContext(ctx, clientMD)
	// transport carrier
	ctx = transport.NewClientContext(ctx, &testTransport{headerCarrier{}})
	// const md
	constMD := metadata.New()
	constMD.Set(constKey, constValue)
	reply, err := Client(WithConstants(constMD))(hs)(ctx, "bar")
	if err != nil {
		t.Fatal(err)
	}
	if reply.(string) != "bar" {
		t.Fatalf("want foo got %v", reply)
	}
}

func TestWithConstants(t *testing.T) {
	md := metadata.Metadata{
		constKey: {constValue},
	}
	options := &options{
		md: metadata.Metadata{
			"override": {"override"},
		},
	}

	WithConstants(md)(options)
	if !reflect.DeepEqual(md, options.md) {
		t.Errorf("want: %v, got: %v", md, options.md)
	}
}

func TestOptions_WithPropagatedPrefix(t *testing.T) {
	options := &options{
		prefix: []string{"override"},
	}
	prefixes := []string{"something", "another"}

	WithPropagatedPrefix(prefixes...)(options)
	if !reflect.DeepEqual(prefixes, options.prefix) {
		t.Error("The prefix must be overridden.")
	}
}

func TestOptions_hasPrefix(t *testing.T) {
	tests := []struct {
		name    string
		options *options
		key     string
		exists  bool
	}{
		{"exists key upper", &options{prefix: []string{"prefix"}}, "PREFIX_true", true},
		{"exists key lower", &options{prefix: []string{"prefix"}}, "prefix_true", true},
		{"not exists key upper", &options{prefix: []string{"prefix"}}, "false_PREFIX", false},
		{"not exists key lower", &options{prefix: []string{"prefix"}}, "false_prefix", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exists := test.options.hasPrefix(test.key)
			if test.exists != exists {
				t.Errorf("key: '%sr', not exists prefixs: %v", test.key, test.options.prefix)
			}
		})
	}
}

func TestEncodeValue(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain ascii", "global-value", "global-value"},
		{"chinese", "张三", "%E5%BC%A0%E4%B8%89"},
		{"emoji", "🙂", "%F0%9F%99%82"},
		{"percent", "100% off", "100%25%20off"},
		{"plus preserved", "+861234567890", "+861234567890"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeValue(tt.in); got != tt.want {
				t.Errorf("encodeValue(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDecodeValue(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain ascii", "global-value", "global-value"},
		{"chinese", "%E5%BC%A0%E4%B8%89", "张三"},
		{"emoji", "%F0%9F%99%82", "🙂"},
		{"percent escaped", "100%25%20off", "100% off"},
		{"plus preserved", "+861234567890", "+861234567890"},
		{"bare percent error fallback", "100% off", "100% off"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeValue(tt.in); got != tt.want {
				t.Errorf("decodeValue(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNonASCIIRoundTripThroughClientServer(t *testing.T) {
	// Simulate a client sending a non-ASCII value through the middleware chain:
	// Client encodes it into the header, then Server reads and decodes it back.
	want := "张三" // any non-ASCII value
	key := "x-md-global-name"

	// Client side: put the value into client metadata, let Client middleware encode it.
	hs := func(ctx context.Context, in any) (any, error) {
		// This handler runs "server-side"; read the value from server context.
		md, ok := metadata.FromServerContext(ctx)
		if !ok {
			return nil, errors.New("no server md")
		}
		if got := md.Get(key); got != want {
			return nil, errors.New("server md value mismatch")
		}
		return in, nil
	}

	// Build the encoded header by running the Client middleware against a carrier.
	clientMD := metadata.New()
	clientMD.Set(key, want)
	clientCtx := metadata.NewClientContext(context.Background(), clientMD)
	carrier := headerCarrier{}
	clientCtx = transport.NewClientContext(clientCtx, &testTransport{carrier})
	_, err := Client()(func(ctx context.Context, in any) (any, error) { return in, nil })(clientCtx, "req")
	if err != nil {
		t.Fatal(err)
	}
	// The header value on the wire must be percent-encoded, not raw UTF-8.
	wire := carrier.Get(key)
	if wire == want {
		t.Fatalf("expected percent-encoded value on the wire, got raw %q", wire)
	}
	if wire != url.PathEscape(want) {
		t.Fatalf("wire value = %q, want %q", wire, url.PathEscape(want))
	}

	// Server side: feed the encoded header into the Server middleware.
	serverCarrier := headerCarrier{}
	serverCarrier.Set(key, wire)
	serverCtx := transport.NewServerContext(context.Background(), &testTransport{serverCarrier})
	_, err = Server()(hs)(serverCtx, "req")
	if err != nil {
		t.Fatal(err)
	}
}
