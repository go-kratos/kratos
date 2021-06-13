package transport

import (
	"strings"
)

// Header is the storage medium used by a Header.
type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

// HeaderCarrier adapts map[string][]string to satisfy the Header interface.
type HeaderCarrier map[string]string

// NewHeaderCarrier new a header carrier.
func NewHeaderCarrier(m map[string][]string) Header {
	hc := make(HeaderCarrier, len(m))
	for k, v := range m {
		if k == "" || len(v) == 0 {
			continue
		}
		hc.Set(k, v[0])
	}
	return hc
}

// Get returns the value associated with the passed key.
func (hc HeaderCarrier) Get(key string) string {
	k := strings.ToLower(key)
	return hc[k]
}

// Set stores the key-value pair.
func (hc HeaderCarrier) Set(key string, value string) {
	k := strings.ToLower(key)
	hc[k] = value
}

// Keys lists the keys stored in this carrier.
func (hc HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}
