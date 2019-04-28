package metadata

import (
	"context"

	"go-common/library/net/metadata"
)

// FromContext get metadata from context.
func FromContext(c context.Context) map[string]string {
	return map[string]string{
		metadata.Color:    metadata.String(c, metadata.Color),
		metadata.Caller:   metadata.String(c, metadata.Caller),
		metadata.Mirror:   metadata.String(c, metadata.Mirror),
		metadata.RemoteIP: metadata.String(c, metadata.RemoteIP),
	}
}

// NewContext new metadata context.
func NewContext(c context.Context, meta map[string]string) context.Context {
	md := metadata.MD{
		metadata.Color:    meta[metadata.Color],
		metadata.Caller:   meta[metadata.Caller],
		metadata.Mirror:   meta[metadata.Mirror],
		metadata.RemoteIP: meta[metadata.RemoteIP],
	}
	return metadata.NewContext(c, md)
}
