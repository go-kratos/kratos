package discovery

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/SeeMusic/kratos/v2/log"
	"github.com/SeeMusic/kratos/v2/registry"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

type mockLogger struct {
	level log.Level
	key   string
	val   string
}

func (l *mockLogger) Log(level log.Level, keyvals ...interface{}) error {
	l.level = level
	l.key = keyvals[0].(string)
	l.val = keyvals[1].(string)
	return nil
}

func TestWithLogger(t *testing.T) {
	b := &builder{}
	WithLogger(&mockLogger{})(b)
}

func TestWithInsecure(t *testing.T) {
	b := &builder{}
	WithInsecure(true)(b)
	if !b.insecure {
		t.Errorf("expected insecure to be true")
	}
}

func TestWithTimeout(t *testing.T) {
	o := &builder{}
	v := time.Duration(123)
	WithTimeout(v)(o)
	if !reflect.DeepEqual(v, o.timeout) {
		t.Errorf("expected %v, got %v", v, o.timeout)
	}
}

type mockDiscovery struct{}

func (m *mockDiscovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (m *mockDiscovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return &testWatch{}, nil
}

func TestBuilder_Scheme(t *testing.T) {
	b := NewBuilder(&mockDiscovery{})
	if !reflect.DeepEqual("discovery", b.Scheme()) {
		t.Errorf("expected %v, got %v", "discovery", b.Scheme())
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

func TestBuilder_Build(t *testing.T) {
	b := NewBuilder(&mockDiscovery{})
	_, err := b.Build(resolver.Target{Scheme: resolver.GetDefaultScheme(), Endpoint: "gprc://authority/endpoint"}, &mockConn{}, resolver.BuildOptions{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
