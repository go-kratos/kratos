package builtin

import (
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
)

var (
	_ metadata.Metadata = &md{}
	_ metadata.Builder  = &Builder{}
)

type metadataKey struct{}

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type md map[string]string

type Builder struct {
}

// Build returns an MD formed by the mapping of key, value ...
// New panics if len(kv) is odd.
func (b *Builder) Build(kvPairs ...string) metadata.Metadata {
	if len(kvPairs)%2 == 1 {
		panic(fmt.Sprintf("metadata: New got the odd number of input pairs for metadata: %d", len(kvPairs)))
	}
	m := md{}
	for i := 0; i < len(kvPairs); i += 2 {
		key := strings.ToLower(kvPairs[i])
		m.Set(key, kvPairs[i+1])
	}
	return m
}

func (m md) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	// attempt to get as is
	val, ok := m[key]
	if ok {
		return val, ok
	}

	return val, ok
}

func (m md) Set(key, val string) {
	key = strings.ToLower(key)
	m[key] = val
}

func (m md) Del(key string) {
	key = strings.ToLower(key)
	// delete key as-is
	delete(m, key)
}

func (m md) Len() int {
	return len(m)
}

// Copy makes a copy of the metadata
func (m md) Copy() metadata.Metadata {
	cmd := md{}
	for k, v := range m {
		cmd[k] = v
	}
	return cmd
}

// Range Traverse all kvs until func return false.
func (m md) Range(f func(k, v string) bool) {
	for key, value := range m {
		if !f(key, value) {
			break
		}
	}
	return
}
