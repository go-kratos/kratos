package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
)

type mockRoundTripper struct{}

func (rt *mockRoundTripper) RoundTrip(_ *http.Request) (resp *http.Response, err error) {
	return
}

type mockCallOption struct {
	needErr bool
}

func (x *mockCallOption) before(_ *callInfo) error {
	if x.needErr {
		return errors.New("option need return err")
	}
	return nil
}

func (x *mockCallOption) after(_ *callInfo, _ *csAttempt) {
	log.Println("run in mockCallOption.after")
}

func TestWithSubset(t *testing.T) {
	co := &clientOptions{}
	o := WithSubset(1)
	o(co)
	if co.subsetSize != 1 {
		t.Error("expected subset size to be 1")
	}
}

func TestWithTransport(t *testing.T) {
	ov := &mockRoundTripper{}
	o := WithTransport(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.transport, ov) {
		t.Errorf("expected transport to be %v, got %v", ov, co.transport)
	}
}

func TestWithTimeout(t *testing.T) {
	ov := 1 * time.Second
	o := WithTimeout(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.timeout, ov) {
		t.Errorf("expected timeout to be %v, got %v", ov, co.timeout)
	}
}

func TestWithBlock(t *testing.T) {
	o := WithBlock()
	co := &clientOptions{}
	o(co)
	if !co.block {
		t.Errorf("expected block to be true, got %v", co.block)
	}
}

func TestWithTLSConfig(t *testing.T) {
	ov := &tls.Config{}
	o := WithTLSConfig(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.tlsConf, ov) {
		t.Errorf("expected tls config to be %v, got %v", ov, co.tlsConf)
	}
}

func TestWithUserAgent(t *testing.T) {
	ov := "kratos"
	o := WithUserAgent(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.userAgent, ov) {
		t.Errorf("expected user agent to be %v, got %v", ov, co.userAgent)
	}
}

func TestWithMiddleware(t *testing.T) {
	o := &clientOptions{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	WithMiddleware(v...)(o)
	if !reflect.DeepEqual(o.middleware, v) {
		t.Errorf("expected middleware to be %v, got %v", v, o.middleware)
	}
}

func TestWithEndpoint(t *testing.T) {
	ov := "some-endpoint"
	o := WithEndpoint(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.endpoint, ov) {
		t.Errorf("expected endpoint to be %v, got %v", ov, co.endpoint)
	}
}

func TestWithRequestEncoder(t *testing.T) {
	o := &clientOptions{}
	v := func(context.Context, string, interface{}) (body []byte, err error) {
		return nil, nil
	}
	WithRequestEncoder(v)(o)
	if o.encoder == nil {
		t.Errorf("expected encoder to be not nil")
	}
}

func TestWithResponseDecoder(t *testing.T) {
	o := &clientOptions{}
	v := func(context.Context, *http.Response, interface{}) error { return nil }
	WithResponseDecoder(v)(o)
	if o.decoder == nil {
		t.Errorf("expected encoder to be not nil")
	}
}

func TestWithErrorDecoder(t *testing.T) {
	o := &clientOptions{}
	v := func(context.Context, *http.Response) error { return nil }
	WithErrorDecoder(v)(o)
	if o.errorDecoder == nil {
		t.Errorf("expected encoder to be not nil")
	}
}

type mockDiscovery struct{}

func (*mockDiscovery) GetService(_ context.Context, _ string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (*mockDiscovery) Watch(_ context.Context, _ string) (registry.Watcher, error) {
	return &mockWatcher{}, nil
}

type mockWatcher struct{}

func (m *mockWatcher) Next() ([]*registry.ServiceInstance, error) {
	instance := &registry.ServiceInstance{
		ID:        "1",
		Name:      "kratos",
		Version:   "v1",
		Metadata:  map[string]string{},
		Endpoints: []string{fmt.Sprintf("http://127.0.0.1:9001?isSecure=%s", strconv.FormatBool(false))},
	}
	time.Sleep(time.Millisecond * 500)
	return []*registry.ServiceInstance{instance}, nil
}

func (*mockWatcher) Stop() error {
	return nil
}

func TestWithDiscovery(t *testing.T) {
	ov := &mockDiscovery{}
	o := WithDiscovery(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.discovery, ov) {
		t.Errorf("expected discovery to be %v, got %v", ov, co.discovery)
	}
}

func TestWithNodeFilter(t *testing.T) {
	ov := func(context.Context, []selector.Node) []selector.Node {
		return []selector.Node{&selector.DefaultNode{}}
	}
	o := WithNodeFilter(ov)
	co := &clientOptions{}
	o(co)
	for _, n := range co.nodeFilters {
		ret := n(context.Background(), nil)
		if len(ret) != 1 {
			t.Errorf("expected node  length to be 1, got %v", len(ret))
		}
	}
}

func TestDefaultRequestEncoder(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBufferString(`{"a":"1", "b": 2}`)))
	r.Header.Set("Content-Type", "application/xml")

	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{"a", 1}
	b, err := DefaultRequestEncoder(context.TODO(), "application/json", v1)
	if err != nil {
		t.Fatal(err)
	}
	v1b := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err = json.Unmarshal(b, v1b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v1b, v1) {
		t.Errorf("expected %v, got %v", v1, v1b)
	}
}

func TestDefaultResponseDecoder(t *testing.T) {
	resp1 := &http.Response{
		Header:     make(http.Header),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{"a":"1", "b": 2}`)),
	}
	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err := DefaultResponseDecoder(context.TODO(), resp1, &v1)
	if err != nil {
		t.Fatal(err)
	}
	if v1.A != "1" {
		t.Errorf("expected %v, got %v", "1", v1.A)
	}
	if v1.B != int64(2) {
		t.Errorf("expected %v, got %v", 2, v1.B)
	}

	resp2 := &http.Response{
		Header:     make(http.Header),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("{badjson}")),
	}
	v2 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err = DefaultResponseDecoder(context.TODO(), resp2, &v2)
	syntaxErr := &json.SyntaxError{}
	if !errors.As(err, &syntaxErr) {
		t.Errorf("expected %v, got %v", syntaxErr, err)
	}
}

func TestDefaultErrorDecoder(t *testing.T) {
	for i := 200; i < 300; i++ {
		resp := &http.Response{Header: make(http.Header), StatusCode: i}
		if DefaultErrorDecoder(context.TODO(), resp) != nil {
			t.Errorf("expected no error, got %v", DefaultErrorDecoder(context.TODO(), resp))
		}
	}
	resp1 := &http.Response{
		Header:     make(http.Header),
		StatusCode: 300,
		Body:       io.NopCloser(bytes.NewBufferString("{\"foo\":\"bar\"}")),
	}
	if DefaultErrorDecoder(context.TODO(), resp1) == nil {
		t.Errorf("expected error, got nil")
	}

	resp2 := &http.Response{
		Header:     make(http.Header),
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewBufferString(`{"code":54321, "message": "hi", "reason": "FOO"}`)),
	}
	err := DefaultErrorDecoder(context.TODO(), resp2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if err.(*kratoserrors.Error).Code != int32(500) {
		t.Errorf("expected %v, got %v", 500, err.(*kratoserrors.Error).Code)
	}
	if err.(*kratoserrors.Error).Message != "hi" {
		t.Errorf("expected %v, got %v", "hi", err.(*kratoserrors.Error).Message)
	}
	if err.(*kratoserrors.Error).Reason != "FOO" {
		t.Errorf("expected %v, got %v", "FOO", err.(*kratoserrors.Error).Reason)
	}
}

func TestCodecForResponse(t *testing.T) {
	resp := &http.Response{Header: make(http.Header)}
	resp.Header.Set("Content-Type", "application/xml")
	c := CodecForResponse(resp)
	if !reflect.DeepEqual("xml", c.Name()) {
		t.Errorf("expected %v, got %v", "xml", c.Name())
	}
}

func TestNewClient(t *testing.T) {
	_, err := NewClient(context.Background(), WithEndpoint("127.0.0.1:8888"))
	if err != nil {
		t.Error(err)
	}
	_, err = NewClient(context.Background(), WithEndpoint("127.0.0.1:9999"), WithTLSConfig(&tls.Config{ServerName: "www.kratos.com", RootCAs: nil}))
	if err != nil {
		t.Error(err)
	}
	_, err = NewClient(context.Background(), WithDiscovery(&mockDiscovery{}), WithEndpoint("discovery:///go-kratos"))
	if err != nil {
		t.Error(err)
	}
	_, err = NewClient(context.Background(), WithDiscovery(&mockDiscovery{}), WithEndpoint("127.0.0.1:8888"))
	if err != nil {
		t.Error(err)
	}
	_, err = NewClient(context.Background(), WithEndpoint("127.0.0.1:8888:xxxxa"))
	if err == nil {
		t.Error("except a parseTarget error")
	}
	_, err = NewClient(context.Background(), WithDiscovery(&mockDiscovery{}), WithEndpoint("https://go-kratos.dev/"))
	if err == nil {
		t.Error("err should not be equal to nil")
	}

	client, err := NewClient(
		context.Background(),
		WithDiscovery(&mockDiscovery{}),
		WithEndpoint("discovery:///go-kratos"),
		WithMiddleware(func(handler middleware.Handler) middleware.Handler {
			t.Logf("handle in middleware")
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				return handler(ctx, req)
			}
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Invoke(context.Background(), http.MethodPost, "/go", map[string]string{"name": "kratos"}, nil, EmptyCallOption{}, &mockCallOption{})
	if err == nil {
		t.Error("err should not be equal to nil")
	}
	err = client.Invoke(context.Background(), http.MethodPost, "/go", map[string]string{"name": "kratos"}, nil, EmptyCallOption{}, &mockCallOption{needErr: true})
	if err == nil {
		t.Error("err should be equal to callOption err")
	}
	client.opts.encoder = func(context.Context, string, interface{}) (body []byte, err error) {
		return nil, errors.New("mock test encoder error")
	}
	err = client.Invoke(context.Background(), http.MethodPost, "/go", map[string]string{"name": "kratos"}, nil, EmptyCallOption{})
	if err == nil {
		t.Error("err should be equal to encoder error")
	}
}
