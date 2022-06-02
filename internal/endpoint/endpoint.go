package endpoint

import (
	"net/url"
	"strconv"
	"strings"
)

// NewEndpoint new an Endpoint URL.
func NewEndpoint(scheme, host string) *url.URL {
	return &url.URL{Scheme: scheme, Host: host}
}

// ParseEndpoint parses an Endpoint URL.
func ParseEndpoint(endpoints []string, scheme string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}
		if u.Scheme == scheme {
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

// Scheme is the scheme of endpoint url.
// examples: scheme="http",isSecure=true get "https"
func Scheme(scheme string, isSecure bool) string {
	if isSecure && scheme != "https" && scheme != "grpcs" {
		scheme += "s"
	}
	return scheme
}
