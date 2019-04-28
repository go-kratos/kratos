package service

import (
	"fmt"
	"math"
	"net"
)

const sixtyTwo uint64 = 62

func encode(num uint64) string {
	var numStr string
	for num != 0 {
		r := num % sixtyTwo
		if r <= 9 {
			numStr = string(byte(r)+48) + numStr
		} else if r <= 35 {
			numStr = string(byte(r)+87) + numStr
		} else {
			numStr = string(byte(r)+29) + numStr
		}
		num = num / sixtyTwo
	}
	return numStr
}

func decode(bs []byte) uint64 {
	f := 0.0
	l := len(bs) - 1
	for _, v := range bs {
		if v >= 97 {
			v -= 87
		} else if v >= 65 {
			v -= 29
		} else {
			v -= 48
		}
		f = f + float64(v)*math.Pow(float64(sixtyTwo), float64(l))
		l = l - 1
	}
	return uint64(f)
}

// ntoIPv6 conver uint32 to ip addr.
func ntoIPv6(sip []string) string {
	if len(sip) != 3 {
		return ""
	}
	ip := make(net.IP, net.IPv6len)
	sum := decode([]byte(sip[0]))
	ip[0] = byte((sum >> 40) & 0xFF)
	ip[1] = byte((sum >> 32) & 0xFF)
	ip[2] = byte((sum >> 24) & 0xFF)
	ip[3] = byte((sum >> 16) & 0xFF)
	ip[4] = byte((sum >> 8) & 0xFF)
	ip[5] = byte(sum & 0xFF)

	sum = decode([]byte(sip[1]))
	ip[6] = byte((sum >> 40) & 0xFF)
	ip[7] = byte((sum >> 32) & 0xFF)
	ip[8] = byte((sum >> 24) & 0xFF)
	ip[9] = byte((sum >> 16) & 0xFF)
	ip[10] = byte((sum >> 8) & 0xFF)
	ip[11] = byte(sum & 0xFF)

	sum = decode([]byte(sip[2]))
	ip[12] = byte((sum >> 24) & 0xFF)
	ip[13] = byte((sum >> 16) & 0xFF)
	ip[14] = byte((sum >> 8) & 0xFF)
	ip[15] = byte(sum & 0xFF)

	return ip.String()
}

// ipv6AtoN conver ip addr to uint32.
func ipv6AtoN(ip net.IP) (sip string) {
	ip = ip.To16()
	if ip == nil {
		return
	}
	sum := uint64(ip[0]) << 40
	sum += uint64(ip[1]) << 32
	sum += uint64(ip[2]) << 24
	sum += uint64(ip[3]) << 16
	sum += uint64(ip[4]) << 8
	sum += uint64(ip[5])
	sip = encode(sum)
	sum = uint64(ip[6]) << 40
	sum += uint64(ip[7]) << 32
	sum += uint64(ip[8]) << 24
	sum += uint64(ip[9]) << 16
	sum += uint64(ip[10]) << 8
	sum += uint64(ip[11])
	sip = sip + ":" + encode(sum)
	sum = uint64(ip[12]) << 24
	sum += uint64(ip[13]) << 16
	sum += uint64(ip[14]) << 8
	sum += uint64(ip[15])
	sip = sip + ":" + encode(sum)
	fmt.Println(sip, "len:", len(sip))
	return
}

// netAtoN conver ipv4 addr to uint32.
func netAtoN(ip net.IP) (sum uint32) {
	ip = ip.To4()
	if ip == nil {
		return
	}
	sum += uint32(ip[0]) << 24
	sum += uint32(ip[1]) << 16
	sum += uint32(ip[2]) << 8
	sum += uint32(ip[3])
	return sum
}

// netNtoA conver uint32 to ipv4 addr.
func netNtoA(sum uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((sum >> 24) & 0xFF)
	ip[1] = byte((sum >> 16) & 0xFF)
	ip[2] = byte((sum >> 8) & 0xFF)
	ip[3] = byte(sum & 0xFF)
	return ip.String()
}

func parseIP(s string) (net.IP, bool) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return net.ParseIP(s), true
		case ':':
			return net.ParseIP(s), false
		}
	}
	return nil, false
}
