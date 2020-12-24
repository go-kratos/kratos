package provider

import "time"

// KeyValue .
// format: json/yaml/text
type KeyValue struct {
	Key       string
	Value     []byte
	Format    string
	Timestamp time.Time
}

// Provider is config provider.
type Provider interface {
	Load() ([]KeyValue, error)
	Watch(key ...string) <-chan KeyValue
}
