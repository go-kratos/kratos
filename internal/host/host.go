package host

import (
	"net"
	"strconv"
)

func isPrivateIP(addr string) bool {
	ip := net.ParseIP(addr)
	if ip4 := ip.To4(); ip4 != nil {
		// Following RFC 4193, Section 3. Local IPv6 Unicast Addresses which says:
		//   The Internet Assigned Numbers Authority (IANA) has reserved the
		//   following three blocks of the IPv4 address space for private internets:
		//     10.0.0.0        -   10.255.255.255  (10/8 prefix)
		//     172.16.0.0      -   172.31.255.255  (172.16/12 prefix)
		//     192.168.0.0     -   192.168.255.255 (192.168/16 prefix)
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	// Following RFC 4193, Section 3. Private Address Space which says:
	//   The Internet Assigned Numbers Authority (IANA) has reserved the
	//   following block of the IPv6 address space for local internets:
	//     FC00::  -  FDFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF (FC00::/7 prefix)
	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}

// Port return a real port.
func Port(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}
	return 0, false
}

// Extract returns a private addr and port.
func Extract(hostPort string, lis net.Listener) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
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
