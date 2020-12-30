package blademaster

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/criticality"
	"github.com/go-kratos/kratos/pkg/net/metadata"

	"github.com/pkg/errors"
)

const (
	// http head
	_httpHeaderUser         = "x1-bmspy-user"
	_httpHeaderTimeout      = "x1-bmspy-timeout"
	_httpHeaderRemoteIP     = "x-backend-bm-real-ip"
	_httpHeaderRemoteIPPort = "x-backend-bm-real-ipport"
)

const (
	_httpHeaderMetadata = "x-bm-metadata-"
)

var _parser = map[string]func(string) interface{}{
	"mirror": func(mirrorStr string) interface{} {
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
	},
	"criticality": func(in string) interface{} {
		if crtl := criticality.Criticality(in); crtl != criticality.EmptyCriticality {
			return string(crtl)
		}
		return string(criticality.Critical)
	},
}

func parseMetadataTo(req *http.Request, to metadata.MD) {
	for rawKey := range req.Header {
		key := strings.ReplaceAll(strings.TrimPrefix(strings.ToLower(rawKey), _httpHeaderMetadata), "-", "_")
		rawValue := req.Header.Get(rawKey)
		var value interface{} = rawValue
		parser, ok := _parser[key]
		if ok {
			value = parser(rawValue)
		}
		to[key] = value
	}
	return
}

func setMetadata(req *http.Request, key string, value interface{}) {
	strV, ok := value.(string)
	if !ok {
		return
	}
	header := fmt.Sprintf("%s%s", _httpHeaderMetadata, strings.ReplaceAll(key, "_", "-"))
	req.Header.Set(header, strV)
}

// setCaller set caller into http request.
func setCaller(req *http.Request) {
	req.Header.Set(_httpHeaderUser, env.AppID)
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
// x-backend-bm-real-ip or X-Real-IP or X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
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
	remote, _, _ = net.SplitHostPort(req.RemoteAddr)
	return
}

func remotePort(req *http.Request) (port string) {
	if port = req.Header.Get(_httpHeaderRemoteIPPort); port != "" && port != "null" {
		return
	}
	return
}
