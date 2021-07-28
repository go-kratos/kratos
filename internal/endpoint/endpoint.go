package endpoint

import (
	"fmt"
	"net/url"
)

func NewEndpoint(scheme, host string, isSecure bool) *url.URL {
	return &url.URL{Scheme: scheme, Host: host, RawQuery: fmt.Sprintf("isSecure=%v", isSecure)}
}

func IsSecure(url *url.URL) bool {
	values, ok := url.Query()["isSecure"]
	if ok && len(values) > 0 && values[0] == "true" {
		return true
	}
	return false
}
