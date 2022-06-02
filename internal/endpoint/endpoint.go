package endpoint

import (
	"net/url"
	"strconv"
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

		// TODO: Compatibility processing
		// Function is to convert grpc:/127.0.0.1/?isSecure=true into grpcs:/127.0.0.1/
		// It will be deleted in about a month
		u = legacyURLToNew(u)
		if u.Scheme == scheme {
			return u.Host, nil
		}
	}
	return "", nil
}

func legacyURLToNew(u *url.URL) *url.URL {
	if u.Scheme == "https" || u.Scheme == "grpcs" {
		return u
	}
	if IsSecure(u) {
		return &url.URL{Scheme: u.Scheme + "s", Host: u.Host}
	}
	return u
}

// IsSecure parses isSecure for Endpoint URL.
// Note: It will be deleted after some time,
// unified use grpcs://127.0.0.1:8080 instead of grpc://127.0.0.1:8080?isSecure=true
func IsSecure(u *url.URL) bool {
	ok, err := strconv.ParseBool(u.Query().Get("isSecure"))
	if err != nil {
		return false
	}
	return ok
}

// Scheme is the scheme of endpoint url.
// examples: scheme="http",isSecure=true get "https"
func Scheme(scheme string, isSecure bool) string {
	if isSecure {
		return scheme + "s"
	}
	return scheme
}
