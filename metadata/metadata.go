package metadata

import (
	"context"
)

type metadataKey struct{}

// Builder is metadata builder
type Builder interface {
	// New returns an MD formed by the mapping of key, value ...
	// New panics if len(kv) is odd.
	Build(kvPairs ...string) Metadata
}

// Metadata is metadata interface
type Metadata interface {
	Get(key string) (val string, ok bool)
	Set(key, val string)
	Del(key string)
	Len() (count int)
	Copy() (md Metadata)
	Range(f func(k, v string) bool)
}

// FromContext returns metadata from the given context
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return nil, ok
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
		return NewContext(context.Background(), patchMd)
	}

	md, ok := FromContext(ctx)
	if !ok {
		return NewContext(ctx, patchMd)
	}

	cmd := md.Copy()
	patchMd.Range(func(k, v string) bool {
		cmd.Set(k, v)
		return true
	})

	return NewContext(ctx, cmd)
}
