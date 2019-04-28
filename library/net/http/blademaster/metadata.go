package blademaster

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	// http head
	_httpHeaderUser         = "x1-bilispy-user"
	_httpHeaderColor        = "x1-bilispy-color"
	_httpHeaderTimeout      = "x1-bilispy-timeout"
	_httpHeaderRemoteIP     = "x-backend-bili-real-ip"
	_httpHeaderRemoteIPPort = "x-backend-bili-real-ipport"
)

// mirror return true if x1-bilispy-mirror in http header and its value is 1 or true.
func mirror(req *http.Request) bool {
	mirrorStr := req.Header.Get("x1-bilispy-mirror")
	if mirrorStr == "" {
		return false
	}
	val, err := strconv.ParseBool(mirrorStr)
	if err != nil {
		log.Warn("blademaster: failed to parse mirror: %+v", errors.Wrap(err, mirrorStr))
		return false
	}
	if !val {
		log.Warn("blademaster: request mirrorStr value :%s is false", mirrorStr)
	}
	return val
}

// setCaller set caller into http request.
func setCaller(req *http.Request) {
	req.Header.Set(_httpHeaderUser, env.AppID)
}

// caller get caller from http request.
func caller(req *http.Request) string {
	return req.Header.Get(_httpHeaderUser)
}

// setColor set color into http request.
func setColor(req *http.Request, color string) {
	req.Header.Set(_httpHeaderColor, color)
}

// color get color from http request.
func color(req *http.Request) string {
	c := req.Header.Get(_httpHeaderColor)
	if c == "" {
		c = env.Color
	}
	return c
}

// setTimeout set timeout into http request.
func setTimeout(req *http.Request, timeout time.Duration) {
	td := int64(timeout / time.Millisecond)
	req.Header.Set(_httpHeaderTimeout, strconv.FormatInt(td, 10))
}

// timeout get timeout from http request.
func timeout(req *http.Request) time.Duration {
	to := req.Header.Get(_httpHeaderTimeout)
	timeout, err := strconv.ParseInt(to, 10, 64)
	if err == nil && timeout > 20 {
		timeout -= 20 // reduce 20ms every time.
	}
	return time.Duration(timeout) * time.Millisecond
}

// remoteIP implements a best effort algorithm to return the real client IP, it parses
// X-BACKEND-BILI-REAL-IP or X-Real-IP or X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func remoteIP(req *http.Request) (remote string) {
	if remote = req.Header.Get(_httpHeaderRemoteIP); remote != "" && remote != "null" {
		return
	}
	var xff = req.Header.Get("X-Forwarded-For")
	if idx := strings.IndexByte(xff, ','); idx > -1 {
		if remote = strings.TrimSpace(xff[:idx]); remote != "" {
			return
		}
	}
	if remote = req.Header.Get("X-Real-IP"); remote != "" {
		return
	}
	remote = req.RemoteAddr[:strings.Index(req.RemoteAddr, ":")]
	return
}

func remotePort(req *http.Request) (port string) {
	if port = req.Header.Get(_httpHeaderRemoteIPPort); port != "" && port != "null" {
		return
	}
	return
}
