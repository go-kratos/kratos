package host

import (
	"net"
	"strconv"
)

var (
	privateAddrs = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10", "fd00::/8"}
)

func isPrivateIP(addr string) bool {
	ipAddr := net.ParseIP(addr)
	for _, privateAddr := range privateAddrs {
		if _, priv, err := net.ParseCIDR(privateAddr); err == nil {
			if priv.Contains(ipAddr) {
				return true
			}
		}
	}
	return false
}

// Port return a real port.
func Port(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}
	return 0, false
}

// Extract returns a private addr and port.
func Extract(hostport string, lis net.Listener) (string, error) {
	addr, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", err
	}
	if lis != nil {
		if p, ok := Port(lis); ok {
			port = strconv.Itoa(p)
		}
	}
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return net.JoinHostPort(addr, port), nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if isPrivateIP(ip.String()) {
				return net.JoinHostPort(ip.String(), port), nil
			}
		}
	}
	return "", nil
}
