package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/vault/api"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx  context.Context
	path string
}

//  WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return Option(func(o *options) {
		o.ctx = ctx
	})
}

// WithPath is config path
func WithPath(p string) Option {
	return Option(func(o *options) {
		o.path = p
	})
}

type source struct {
	client  *api.Client
	options *options
}

func New(client *api.Client, opts ...Option) (config.Source, error) {
	options := &options{
		ctx:  context.Background(),
		path: "",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.path == "" {
		return nil, errors.New("path invalid")
	}

	return &source{
		client:  client,
		options: options,
	}, nil
}

func extractKV(secretData map[string]interface{}) []*config.KeyValue {
	kvs := make([]*config.KeyValue, 0)
	for key, item := range secretData {
		switch v := item.(type) {
		case []byte:
			kvs = append(kvs, &config.KeyValue{Key: key, Value: v})
		case string:
			kvs = append(kvs, &config.KeyValue{Key: key, Value: []byte(v)})
		case map[string]interface{}:
			for _, subItem := range extractKV(v) {
				kvs = append(kvs, &config.KeyValue{Key: fmt.Sprintf("%s/%s", key, subItem.Key), Value: subItem.Value})
			}
		default:
			kvs = append(kvs, &config.KeyValue{Key: key, Value: []byte(fmt.Sprint(v))})
		}
	}
	return kvs
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	secret, err := s.client.Logical().Read(s.options.path)
	if err != nil {
		return nil, err
	}

	return extractKV(secret.Data), nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s)
}
