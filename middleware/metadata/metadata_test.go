package metadata

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string { return http.Header(hc).Get(key) }

func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }

// Keys lists the keys stored in this carrier.
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

type testTransport struct{ header headerCarrier }

func (tr *testTransport) Kind() transport.Kind            { return transport.KindHTTP }
func (tr *testTransport) Endpoint() string                { return "" }
func (tr *testTransport) Operation() string               { return "" }
func (tr *testTransport) RequestHeader() transport.Header { return tr.header }
func (tr *testTransport) ReplyHeader() transport.Header   { return tr.header }

func TestSever(t *testing.T) {
	var (
		globalKey   = "x-md-global-key"
		globalValue = "global-value"
		localKey    = "x-md-local-key"
		localValue  = "local-value"
		constKey    = "x-md-local-const"
		constValue  = "x-md-local-const"
	)
	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
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
	hs := func(ctx context.Context, in interface{}) (interface{}, error) {
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
