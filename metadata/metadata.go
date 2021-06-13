package metadata

import (
	"strings"
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

// New creates an MD from a given key-values map.
func New(mds ...map[string][]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, v := range m {
			if k == "" {
				continue
			}
			key := strings.ToLower(k)
			if len(v) > 0 && v[0] != "" {
				md[key] = v[0]
			}
		}
	}
	return md
}

// Get returns the value associated with the passed key.
func (m Metadata) Get(key string) string {
	k := strings.ToLower(key)
	return m[k]
}

// Set stores the key-value pair.
func (m Metadata) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	k := strings.ToLower(key)
	m[k] = value
}

// Clone returns a deep copy of Metadata
func (m Metadata) Clone() Metadata {
	md := Metadata{}
	for k, v := range m {
		md[k] = v
	}
	return md
}
