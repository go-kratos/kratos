package memory

import (
	"context"

	"github.com/go-kratos/kratos/v2/config"
)

type keyValueKey struct{}

func withData(d []byte, f string) config.Option {
	return func(o *config.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, keyValueKey{}, &config.KeyValue{
			Key:      f,
			Value:    d,
			Metadata: nil,
		})
	}
}

// WithKeyValue allows a keyvalue to be set
func WithKeyValue(cs *config.KeyValue) config.Option {
	return func(o *config.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, keyValueKey{}, cs)
	}
}

// WithJSON allows the source data to be set to json
func WithJSON(d []byte) config.Option {
	return withData(d, "json")
}

// WithYAML allows the source data to be set to yaml
func WithYAML(d []byte) config.Option {
	return withData(d, "yaml")
}
