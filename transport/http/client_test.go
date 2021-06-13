package http

import (
	"context"
	"github.com/go-kratos/kratos/v2/registry"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-kratos/kratos/v2/encoding/xml"
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
	ov := 1*time.Second
	o := WithTimeout(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.timeout, ov)
}

func TestWithBalancer(t *testing.T) {

}

func TestWithUserAgent(t *testing.T) {
	ov := "kratos"
	o := WithUserAgent(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.userAgent, ov)
}


func TestWithMiddleware(t *testing.T) {
}

func TestWithEndpoint(t *testing.T) {
	ov := "some-endpoint"
	o := WithEndpoint(ov)
	co := &clientOptions{}
	o(co)
	assert.Equal(t, co.endpoint, ov)
}

func TestWithRequestEncoder(t *testing.T) {

}

func TestWithResponseDecoder(t *testing.T) {

}

func TestWithErrorDecoder(t *testing.T) {
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

func TestCodecForResponse(t *testing.T) {
	resp := &nethttp.Response{Header:make(nethttp.Header)}
	resp.Header.Set("Content-Type", "application/xml")
	c := CodecForResponse(resp)
	assert.Equal(t, "xml", c.Name())
}
