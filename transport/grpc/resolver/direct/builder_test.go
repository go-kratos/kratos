package direct

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"testing"
)

func TestDirectBuilder_Scheme(t *testing.T) {
	b := NewBuilder()
	assert.Equal(t, "direct", b.Scheme())
}

type mockConn struct {
}

func (m *mockConn) UpdateState(resolver.State) error {
	return nil
}

func (m *mockConn) ReportError(error) {}

func (m *mockConn) NewAddress(addresses []resolver.Address) {}

func (m *mockConn) NewServiceConfig(serviceConfig string) {}

func (m *mockConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}

func TestDirectBuilder_Build(t *testing.T) {
	b := NewBuilder()
	r, err := b.Build(resolver.Target{}, &mockConn{}, resolver.BuildOptions{})
	assert.NoError(t, err)
	r.ResolveNow(resolver.ResolveNowOptions{})
}
