package provider

import "time"

// KeyValue is config key value.
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
	Watch() <-chan KeyValue
}
