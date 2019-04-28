package model

import (
	"fmt"
)

// LoginLog login log.
type LoginLog struct {
	Mid       int64  `json:"mid"`
	Timestamp int64  `json:"timestamp"`
	LoginIP   int64  `json:"loginip"`
	Type      int64  `json:"type"`
	Server    string `json:"server"`
}

// LoginLogResp login log.
type LoginLogResp struct {
	Mid       int64  `json:"mid"`
	Timestamp int64  `json:"timestamp"`
	LoginIP   string `json:"loginip"`
	Type      int64  `json:"type"`
	Server    string `json:"server"`
}

// Format format login log to login log resp.
func Format(l *LoginLog) *LoginLogResp {
	if l == nil {
		return nil
	}
	return &LoginLogResp{
		Mid:       l.Mid,
		Timestamp: l.Timestamp,
		LoginIP:   InetNtoA(l.LoginIP),
		Type:      l.Type,
		Server:    l.Server,
	}
}

// InetNtoA .
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
