package metadata

import (
	"context"
	"fmt"
	"strings"
)

type metadataKey struct{}
type mdOutgoingKey struct{}

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata struct {
	md map[string]string
}

// New returns an MD formed by the mapping of key, value ...
// New panics if len(kv) is odd.
func New(kv ...string) Metadata {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: New got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := Metadata{md: map[string]string{}}
	for i := 0; i < len(kv); i += 2 {
		key := strings.ToLower(kv[i])
		md.Set(key, kv[i+1])
	}
	return md
}

func (m Metadata) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	// attempt to get as is
	val, ok := m.md[key]
	if ok {
		return val, ok
	}

	return val, ok
}

func (m Metadata) Set(key, val string) {
	key = strings.ToLower(key)
	m.md[key] = val
}

func (m Metadata) Del(key string) {
	key = strings.ToLower(key)
	// delete key as-is
	delete(m.md, key)
}

func (m Metadata) Len() int {
	return len(m.md)
}

// Copy makes a copy of the metadata
func (m Metadata) Copy() Metadata {
	cmd := Metadata{md: make(map[string]string, m.Len())}
	for k, v := range m.md {
		cmd.md[k] = v
	}
	return cmd
}

func (m Metadata) Merge(patchMd Metadata) Metadata {
	cmd := m.Copy()
	for k, v := range patchMd.md {
		if v != "" {
			cmd.md[k] = v
		}
	}
	return cmd
}

// Keys lists the keys stored in this carrier.
func (m Metadata) Keys() []string {
	keys := make([]string, 0, m.Len())
	for key := range m.md {
		keys = append(keys, key)
	}
	return keys
}

// Range Traverse all kvs until func return false.
func (m Metadata) Range(f func(k, v string) bool) {
	for key, value := range m.md {
		if !f(key, value) {
			break
		}
	}
	return
}

// FromContext returns metadata from the given context
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return Metadata{md: map[string]string{}}, ok
	}
	return md, ok
}

// NewContext creates a new context with the given metadata
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataKey{}, md)
}

// MergeContext merges metadata to existing metadata, overwriting if specified
func MergeContext(ctx context.Context, patchMd Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md, _ := FromContext(ctx)
	cmd := md.Merge(patchMd)
	return NewContext(ctx, cmd)
}

// FromOutgoingContext returns metadata from the given context.
func FromOutgoingContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(mdOutgoingKey{}).(Metadata)
	if !ok {
		return Metadata{md: map[string]string{}}, ok
	}
	return md, ok
}

// NewOutgoingContext creates a new context with outgoing md attached. If used
// in conjunction with AppendToOutgoingContext, NewOutgoingContext will
// overwrite any previously-appended metadata.
func NewOutgoingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, md)
}

// MergeContext merges metadata to existing metadata, overwriting if specified
func MergeOutgoingContext(ctx context.Context, patchMd Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md, _ := FromOutgoingContext(ctx)
	cmd := md.Merge(patchMd)
	return NewOutgoingContext(ctx, cmd)
}
