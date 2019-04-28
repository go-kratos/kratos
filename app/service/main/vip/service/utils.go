package service

import (
	"net"
	"time"
)

// EndOfDay end of day
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// IPStr get ip string.
func IPStr(ip net.IP) string {
	if ip == nil {
		return ""
	}
	return ip.String()
}
