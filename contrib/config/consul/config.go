package consul

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx  context.Context
	path string
}

// WithContext with registry context.
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

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	kv, _, err := s.client.KV().List(s.options.path, nil)
	if err != nil {
		return nil, err
	}

	pathPrefix := s.options.path
	if !strings.HasSuffix(s.options.path, "/") {
		pathPrefix = pathPrefix + "/"
	}
	kvs := make([]*config.KeyValue, 0)
	for _, item := range kv {
		k := strings.TrimPrefix(item.Key, pathPrefix)
		if k == "" {
			continue
		}
		kvs = append(kvs, &config.KeyValue{
			Key:    k,
			Value:  item.Value,
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		})
	}
	return kvs, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s)
}
