package transport

import (
	"strings"
)

// HeaderCarrier adapts map[string][]string to satisfy the Header interface.
type HeaderCarrier map[string][]string

// Get returns the value associated with the passed key.
func (hc HeaderCarrier) Get(key string) string {
	k := strings.ToLower(key)
	v := hc[k]
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

// Set stores the key-value pair.
func (hc HeaderCarrier) Set(key string, value string) {
	k := strings.ToLower(key)
	hc[k] = []string{value}
}

// Keys lists the keys stored in this carrier.
func (hc HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}
