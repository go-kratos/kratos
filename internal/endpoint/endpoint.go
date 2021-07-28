package endpoint

import (
	"net/url"
	"strconv"
)

func NewEndpoint(scheme, host string, isSecure bool) *url.URL {
	var query string
	if isSecure {
		query = "isSecure=true"
	}
	return &url.URL{Scheme: scheme, Host: host, RawQuery: query}
}

func IsSecure(url *url.URL) bool {
	ok, err := strconv.ParseBool(url.Query().Get("isSecure"))
	if err != nil {
		return false
	}
	return ok
}
