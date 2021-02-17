package host

import (
	"net"
	"testing"
)

func TestPrivateIP(t *testing.T) {
	tests := []struct {
		addr   string
		expect bool
	}{
		{"10.1.0.1", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"8.8.8.8", false},
	}
	for _, test := range tests {
		t.Run(test.addr, func(t *testing.T) {
			res := isPrivateIP(test.addr)
			if res != test.expect {
				t.Fatalf("expected %t got %t", test.expect, res)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		addr   string
		expect string
	}{
		{"127.0.0.1:80", "127.0.0.1:80"},
		{"10.0.0.1:80", "10.0.0.1:80"},
		{"172.16.0.1:80", "172.16.0.1:80"},
		{"192.168.1.1:80", "192.168.1.1:80"},
		{"0.0.0.0:80", ""},
		{"[::]:80", ""},
		{":80", ""},
	}
	for _, test := range tests {
		t.Run(test.addr, func(t *testing.T) {
			res, err := Extract(test.addr, nil)
			if err != nil {
				t.Fatal(err)
			}
			if res != test.expect && (test.expect == "" && test.addr == test.expect) {
				t.Fatalf("expected %s got %s", test.expect, res)
			}
		})
	}
}

func TestPort(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	port, ok := Port(lis)
	if !ok || port == 0 {
		t.Fatalf("expected: %s got %d", lis.Addr().String(), port)
	}
}
