package tracing

import (
	"google.golang.org/grpc/metadata"
)

// MetadataCarrier is grpc metadata carrier
type MetadataCarrier metadata.MD

// Get returns the value associated with the passed key.
func (mc MetadataCarrier) Get(key string) string {
	values := metadata.MD(mc).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// Set stores the key-value pair.
func (mc MetadataCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc MetadataCarrier) Keys() []string {
	keys := make([]string, 0, metadata.MD(mc).Len())
	for key := range metadata.MD(mc) {
		keys = append(keys, key)
	}
	return keys
}

// Del delete key
func (mc MetadataCarrier) Del(key string) {
	delete(mc, key)
}

// Clone copy MetadataCarrier
func (mc MetadataCarrier) Clone() MetadataCarrier {
	return MetadataCarrier(metadata.MD(mc).Copy())
}
