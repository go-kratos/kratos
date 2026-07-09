package grpc

import (
	"testing"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

// nopSubConn embeds the interface only to satisfy it; Build never calls its
// methods, it just uses the SubConn as a map key and reads the address.
type nopSubConn struct{ balancer.SubConn }

// Regression for #3842: init() captures selector.GlobalSelector() before
// client.go's init() sets the default, so the registered builder can be nil.
// Build must fall back to the global selector instead of panicking.
func TestBalancerBuildNilBuilder(t *testing.T) {
	b := &balancerBuilder{builder: nil}
	info := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			nopSubConn{}: {Address: resolver.Address{Addr: "127.0.0.1:9000"}},
		},
	}
	if p := b.Build(info); p == nil {
		t.Fatal("expected a non-nil picker")
	}
}
