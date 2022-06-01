package endpoint

import (
	"net/url"
	"strconv"
	"strings"
)

// NewEndpoint new an Endpoint URL.
func NewEndpoint(scheme, host string, isSecure bool) *url.URL {
	if isSecure && !strings.HasSuffix(scheme, "s") {
		scheme += "s"
	}
	return &url.URL{Scheme: scheme, Host: host}
}

// ParseEndpoint parses an Endpoint URL.
func ParseEndpoint(endpoints []string, scheme string, isSecure bool) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}
		if strings.TrimSuffix(u.Scheme, "s") == strings.TrimSuffix(scheme, "s") &&
			IsSecure(u) == isSecure {
			return u.Host, nil
		}
	}
	return "", nil
}

// IsSecure parses isSecure for Endpoint URL.
// Note: It will be deleted after some time,
// unified use grpcs://127.0.0.1:8080 instead of grpc://127.0.0.1:8080?isSecure=true
func IsSecure(u *url.URL) bool {
	if strings.HasSuffix(u.Scheme, "s") {
		return true
	}
	ok, err := strconv.ParseBool(u.Query().Get("isSecure"))
	if err != nil {
		return false
	}
	return ok
}
