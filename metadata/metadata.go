package metadata

import (
	"context"
	"strings"
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

// New creates an MD from a given key-values map.
func New(m map[string][]string) Metadata {
	md := Metadata{}
	for k, v := range m {
		key := strings.ToLower(k)
		if len(v) > 0 {
			md[key] = v[0]
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
	k := strings.ToLower(key)
	m[k] = value
}

// Keys lists the keys stored in this carrier.
func (m Metadata) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Pairs returns all metadata to key/value pairs.
func (m Metadata) Pairs() []string {
	var kvs = make([]string, len(m)*2)
	for k, v := range m {
		kvs = append(kvs, k, v)
	}
	return kvs
}

type mdIncomingKey struct{}
type mdOutgoingKey struct{}

// NewIncomingContext creates a new context with incoming md attached.
func NewIncomingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

// NewOutgoingContext creates a new context with outgoing md attached.
func NewOutgoingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, md)
}

// FromIncomingContext returns the incoming metadata in ctx if it exists.
func FromIncomingContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(mdIncomingKey{}).(Metadata)
	return md, ok
}

// FromOutgoingContext returns the outgoing metadata in ctx if it exists.
func FromOutgoingContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(mdOutgoingKey{}).(Metadata)
	return md, ok
}
