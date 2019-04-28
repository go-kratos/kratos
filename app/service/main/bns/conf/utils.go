package conf

import (
	"fmt"
	"net"
	"strings"
)

// ProtoAddr ProtoAddr
type ProtoAddr struct {
	Proto, Net, Addr string
}

func (p ProtoAddr) String() string {
	return p.Proto + "://" + p.Addr
}

// socketPath tests if a given address describes a domain socket,
// and returns the relevant path part of the string if it is.
func socketPath(addr string) string {
	if !strings.HasPrefix(addr, "unix://") {
		return ""
	}
	return strings.TrimPrefix(addr, "unix://")
}

// ClientListener is used to format a listener for a
// port on a addr, whatever HTTP, HTTPS, DNS or RPC.
func (c *Config) ClientListener(addr string, port int) (net.Addr, error) {
	if path := socketPath(addr); path != "" {
		return &net.UnixAddr{Name: path, Net: "unix"}, nil
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, fmt.Errorf("Failed to parse IP: %v", addr)
	}
	return &net.TCPAddr{IP: ip, Port: port}, nil
}

// DNSAddrs returns the bind addresses for the DNS server.
func (c *Config) DNSAddrs() ([]ProtoAddr, error) {
	if c.DNS == nil {
		return nil, nil
	}
	a, err := c.ClientListener(c.DNS.Addr, c.DNS.Port)
	if err != nil {
		return nil, err
	}

	addrs := []ProtoAddr{
		{"dns", "tcp", a.String()},
		{"dns", "udp", a.String()},
	}
	return addrs, nil
}

// HTTPAddrs returns the bind addresses for the HTTP server and
// the application protocol which should be served, e.g. 'http'
// or 'https'.
func (c *Config) HTTPAddrs() ([]ProtoAddr, error) {
	var addrs []ProtoAddr

	if c.HTTP != nil {
		a, err := c.ClientListener(c.HTTP.Addr, c.HTTP.Port)
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, ProtoAddr{"http", a.Network(), a.String()})
	}
	return addrs, nil
}
