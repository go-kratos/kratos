package transport

// Request is rpc request
type Request struct {
	// Method value contains: [ POST|GET|PUT|DELETE|PATCH|HEAD|TRACE|OPTIONS|CONNECT ]
	Method string
	// FullPath is full path
	FullPath string
	// Path is path pattern template
	PathPattern string
	// Metadata is http header or grpc metadata
	Metadata Metadata
	// RemoteAddr is peer endpoint address
	RemoteAddr string
	// Query is request query string
	Query string
}

// Metadata is request metadata
type Metadata interface {
	// Get returns the value associated with the passed key.
	Get(key string) string
	// Set stores the key-value pair.
	Set(key string, value string)
	// Keys lists the keys stored in this metadata.
	Keys() []string
	// Del delete key
	Del(key string)
	// Clone returns a copy of metadata or nil if metadata is nil.
	Clone() Metadata
}
