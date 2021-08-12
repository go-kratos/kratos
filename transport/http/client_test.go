package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
}

func (rt *mockRoundTripper) RoundTrip(req *nethttp.Request) (resp *nethttp.Response, err error) {
	return
}

func TestWithTransport(t *testing.T) {
	ov := &mockRoundTripper{}
	o := WithTransport(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.transport, ov)
}

func TestWithTimeout(t *testing.T) {
	ov := 1 * time.Second
	o := WithTimeout(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.timeout, ov)
}

func TestWithBlock(t *testing.T) {
	o := WithBlock()
	co := &clientOptions{}
	o(co)
	assert.True(t, co.block)
}

func TestWithBalancer(t *testing.T) {

}

func TestWithTLSConfig(t *testing.T) {
	ov := &tls.Config{}
	o := WithTLSConfig(ov)
	co := &clientOptions{}
	o(co)
	assert.Same(t, ov, co.tlsConf)
}

func TestWithUserAgent(t *testing.T) {
	ov := "kratos"
	o := WithUserAgent(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.userAgent, ov)
}

func TestWithMiddleware(t *testing.T) {
	o := &clientOptions{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	WithMiddleware(v...)(o)
	assert.Equal(t, v, o.middleware)
}

func TestWithEndpoint(t *testing.T) {
	ov := "some-endpoint"
	o := WithEndpoint(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.endpoint, ov)
}

func TestWithRequestEncoder(t *testing.T) {
	o := &clientOptions{}
	v := func(ctx context.Context, contentType string, in interface{}) (body []byte, err error) {
		return nil, nil
	}
	WithRequestEncoder(v)(o)
	assert.NotNil(t, o.encoder)
}

func TestWithResponseDecoder(t *testing.T) {
	o := &clientOptions{}
	v := func(ctx context.Context, res *nethttp.Response, out interface{}) error { return nil }
	WithResponseDecoder(v)(o)
	assert.NotNil(t, o.decoder)
}

func TestWithErrorDecoder(t *testing.T) {
	o := &clientOptions{}
	v := func(ctx context.Context, res *nethttp.Response) error { return nil }
	WithErrorDecoder(v)(o)
	assert.NotNil(t, o.errorDecoder)
}

type mockDiscovery struct {
}

func (*mockDiscovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (*mockDiscovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return nil, nil
}

func TestWithDiscovery(t *testing.T) {
	ov := &mockDiscovery{}
	o := WithDiscovery(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.discovery, ov)
}

func TestDefaultRequestEncoder(t *testing.T) {
	req1 := &nethttp.Request{
		Header: make(nethttp.Header),
		Body:   ioutil.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	req1.Header.Set("Content-Type", "application/xml")

	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{"a", 1}
	b, err1 := DefaultRequestEncoder(context.TODO(), "application/json", v1)
	assert.Nil(t, err1)
	v1b := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err1 = json.Unmarshal(b, v1b)
	assert.Nil(t, err1)
	assert.Equal(t, v1, v1b)
}

func TestDefaultResponseDecoder(t *testing.T) {
	resp1 := &nethttp.Response{
		Header:     make(nethttp.Header),
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"a\":\"1\", \"b\": 2}")),
	}
	v1 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err1 := DefaultResponseDecoder(context.TODO(), resp1, &v1)
	assert.Nil(t, err1)
	assert.Equal(t, "1", v1.A)
	assert.Equal(t, int64(2), v1.B)

	resp2 := &nethttp.Response{
		Header:     make(nethttp.Header),
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("{badjson}")),
	}
	v2 := &struct {
		A string `json:"a"`
		B int64  `json:"b"`
	}{}
	err2 := DefaultResponseDecoder(context.TODO(), resp2, &v2)
	terr1 := &json.SyntaxError{}
	assert.ErrorAs(t, err2, &terr1)
}

func TestDefaultErrorDecoder(t *testing.T) {
	for i := 200; i < 300; i++ {
		resp := &nethttp.Response{Header: make(nethttp.Header), StatusCode: i}
		assert.Nil(t, DefaultErrorDecoder(context.TODO(), resp))
	}
	resp1 := &nethttp.Response{
		Header:     make(nethttp.Header),
		StatusCode: 300,
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"foo\":\"bar\"}")),
	}
	assert.Error(t, DefaultErrorDecoder(context.TODO(), resp1))

	resp2 := &nethttp.Response{
		Header:     make(nethttp.Header),
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"code\":54321, \"message\": \"hi\", \"reason\": \"FOO\"}")),
	}
	err2 := DefaultErrorDecoder(context.TODO(), resp2)
	assert.Error(t, err2)
	assert.Equal(t, int32(500), err2.(*errors.Error).GetCode())
	assert.Equal(t, "hi", err2.(*errors.Error).GetMessage())
	assert.Equal(t, "FOO", err2.(*errors.Error).GetReason())
}

func TestCodecForResponse(t *testing.T) {
	resp := &nethttp.Response{Header: make(nethttp.Header)}
	resp.Header.Set("Content-Type", "application/xml")
	c := CodecForResponse(resp)
	assert.Equal(t, "xml", c.Name())
}
