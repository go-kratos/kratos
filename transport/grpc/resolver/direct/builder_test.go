package direct

import (
	"errors"
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

type mockConn struct {
	needUpdateStateErr bool
}

func (m *mockConn) UpdateState(resolver.State) error {
	if m.needUpdateStateErr {
		return errors.New("mock test needUpdateStateErr")
	}
	return nil
}

func (m *mockConn) ReportError(error) {}

func (m *mockConn) NewAddress(_ []resolver.Address) {}

func (m *mockConn) NewServiceConfig(_ string) {}

func (m *mockConn) ParseServiceConfig(_ string) *serviceconfig.ParseResult {
	return nil
}

func TestDirectBuilder_Build(t *testing.T) {
	b := NewBuilder()
	r, err := b.Build(resolver.Target{}, &mockConn{}, resolver.BuildOptions{})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	r.Close()

	// need update state err
	_, err = b.Build(resolver.Target{}, &mockConn{needUpdateStateErr: true}, resolver.BuildOptions{})
	if err == nil {
		t.Errorf("expect needUpdateStateErr, got nil")
	}
}
