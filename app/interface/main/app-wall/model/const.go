package model

import (
	"crypto/md5"
	"encoding/hex"
	"net"
)

const (
	TypeIOS         = "ios"
	TypeAndriod     = "android"
	GdtIOSAppID     = "736536022"
	GdtAndroidAppID = "100951776"

	ChannelToutiao = "toutiao"
	ChannelShike   = "2883"
	ChannelDontin  = "415209141"
)

type GdtKey struct {
	Encrypt string
	Sign    string
}

var (
	ChannelGdt = map[string]*GdtKey{
		"1439767": &GdtKey{Encrypt: "BAAAAAAAAAAAFfgX", Sign: "ee358e8dccbbc4ba"},
		"406965":  &GdtKey{Encrypt: "BAAAAAAAAAAABjW1", Sign: "a45cbd2d4c5344b3"},
		"7799673": &GdtKey{Encrypt: "BAAAAAAAAAAAdwN5", Sign: "54b6deffcd64b6b0"},
	}

	AppIDGdt = map[string]string{
		TypeIOS:     GdtIOSAppID,
		TypeAndriod: GdtAndroidAppID,
	}
)

func GdtIMEI(imei string) (gdtImei string) {
	if imei == "" {
		return
	}
	bs := md5.Sum([]byte(imei))
	gdtImei = hex.EncodeToString(bs[:])
	return
}

// InetAtoN conver ip addr to uint32.
func InetAtoN(s string) (sum uint32) {
	ip := net.ParseIP(s)
	if ip == nil {
		return
	}
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

// InetNtoA conver uint32 to ip addr.
func InetNtoA(sum uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((sum >> 24) & 0xFF)
	ip[1] = byte((sum >> 16) & 0xFF)
	ip[2] = byte((sum >> 8) & 0xFF)
	ip[3] = byte(sum & 0xFF)
	return ip.String()
}

// IsIPv4 is ipv4
func IsIPv4(addr string) bool {
	ipv := net.ParseIP(addr)
	if ip := ipv.To4(); ip != nil {
		return true
	} else {
		return false
	}
}
