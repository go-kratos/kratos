package direct

import (
	"reflect"
	"testing"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestDirectBuilder_Scheme(t *testing.T) {
	b := NewBuilder()
	if !reflect.DeepEqual(b.Scheme(), "direct") {
		t.Errorf("expect %v, got %v", "direct", b.Scheme())
	}
}

type mockConn struct{}

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
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
}
