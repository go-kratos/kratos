package http

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/transport"
)

var _ transport.Metadata = HeaderCarrier{}

// HeaderCarrier adapts http.Header to satisfy the TextMapCarrier interface.
type HeaderCarrier http.Header

// Get returns the value associated with the passed key.
func (hc HeaderCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

// Set stores the key-value pair.
func (hc HeaderCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}

// Del delete key
func (hc HeaderCarrier) Del(key string) {
	http.Header(hc).Del(key)
}

// Clone copy HeaderCarrier
func (hc HeaderCarrier) Clone() transport.Metadata {
	return HeaderCarrier(http.Header(hc).Clone())
}
