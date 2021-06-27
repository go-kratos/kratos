package host

import (
	"net"
	"testing"
)

func TestValidIP(t *testing.T) {
	tests := []struct {
		addr   string
		expect bool
	}{
		{"127.0.0.1", false},
		{"255.255.255.255", false},
		{"0.0.0.0", false},
		{"localhost", false},
		{"10.1.0.1", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"8.8.8.8", true},
		{"1.1.1.1", true},
		{"9.255.255.255", true},
		{"10.0.0.0", true},
		{"10.255.255.255", true},
		{"11.0.0.0", true},
		{"172.15.255.255", true},
		{"172.16.0.0", true},
		{"172.16.255.255", true},
		{"172.23.18.255", true},
		{"172.31.255.255", true},
		{"172.31.0.0", true},
		{"172.32.0.0", true},
		{"192.167.255.255", true},
		{"192.168.0.0", true},
		{"192.168.255.255", true},
		{"192.169.0.0", true},
		{"fbff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", true},
		{"fc00::", true},
		{"fcff:1200:0:44::", true},
		{"fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", true},
		{"fe00::", true},
	}
	for _, test := range tests {
		t.Run(test.addr, func(t *testing.T) {
			res := isValidIP(test.addr)
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

func TestExtractHostPort(t *testing.T) {
	host, port, err := ExtractHostPort("127.0.0.1:8000")
	if err != nil {
		t.Fatalf("expected: %v got %v", nil, err)
	}
	t.Logf("host port: %s,  %d", host, port)

	host, port, err = ExtractHostPort("www.bilibili.com:80")
	if err != nil {
		t.Fatalf("expected: %v got %v", nil, err)
	}
	t.Logf("host port: %s,  %d", host, port)

	host, port, err = ExtractHostPort("consul://2/33")
	if err == nil {
		t.Fatalf("expected: not nil got %v", nil)
	}
}
