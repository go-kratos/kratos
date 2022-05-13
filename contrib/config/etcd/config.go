package etcd

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx    context.Context
	paths  []string
	prefix bool
}

// WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithPath is config path
func WithPath(paths ...string) Option {
	return func(o *options) {
		o.paths = paths
	}
}

// WithPrefix is config prefix
func WithPrefix(prefix bool) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

type source struct {
	client  *clientv3.Client
	options *options
}

func New(client *clientv3.Client, opts ...Option) (config.Source, error) {
	options := &options{
		ctx:    context.Background(),
		prefix: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	if len(options.paths) == 0 {
		return nil, errors.New("paths can't empty")
	}

	return &source{
		client:  client,
		options: options,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}

	kvs := make([]*config.KeyValue, 0)
	for _, path := range s.options.paths {
		rsp, err := s.client.Get(s.options.ctx, path, opts...)
		if err != nil {
			return nil, err
		}

		for _, item := range rsp.Kvs {
			k := string(item.Key)
			ext := strings.TrimPrefix(filepath.Ext(k), ".")

			kvs = append(kvs, &config.KeyValue{
				Key:    k,
				Value:  item.Value,
				Format: ext,
			})
		}
	}
	return kvs, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s), nil
}
