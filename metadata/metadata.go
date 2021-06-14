package metadata

import (
	"context"
	"fmt"
	"strings"
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

// New creates an MD from a given key-values map.
func New(mds ...map[string]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, v := range m {
			if k == "" {
				continue
			}
			key := strings.ToLower(k)
			if len(v) > 0 && v != "" {
				md[key] = v
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

// Range iterate over element in metadata.
func (m Metadata) Range(f func(k, v string) bool) {
	for k, v := range m {
		ret := f(k, v)
		if ret == false {
			break
		}
	}
}

// Clone returns a deep copy of Metadata
func (m Metadata) Clone() Metadata {
	md := Metadata{}
	for k, v := range m {
		md[k] = v
	}
	return md
}

type serverMetadataKey struct{}

// NewServerContext creates a new context with client md attached.
func NewServerContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// FromServerContext returns the server metadata in ctx if it exists.
func FromServerContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(serverMetadataKey{}).(Metadata)
	return md, ok
}

type clientMetadataKey struct{}

// NewClientContext creates a new context with client md attached.
func NewClientContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, clientMetadataKey{}, md)
}

// FromClientContext returns the client metadata in ctx if it exists.
func FromClientContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(clientMetadataKey{}).(Metadata)
	return md, ok
}

// AppendToClientContext returns a new context with the provided kv merged
// with any existing metadata in the context.
func AppendToClientContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToOutgoingContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for i := 0; i < len(kv); i += 2 {
		md.Set(kv[i], kv[i+1])
	}
	return NewClientContext(ctx, md)
}

// MergeToClientContext merge new metadata into ctx.
func MergeToClientContext(ctx context.Context, cmd Metadata) context.Context {
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for k, v := range cmd {
		md[k] = v
	}
	return NewClientContext(ctx, md)
}
