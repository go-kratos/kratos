package etcd

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/config"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx    context.Context
	path   string
	prefix bool
}

// WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithPath is config path
func WithPath(p string) Option {
	return func(o *options) {
		o.path = p
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
		path:   "",
		prefix: false,
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
	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}

	rsp, err := s.client.Get(s.options.ctx, s.options.path, opts...)
	if err != nil {
		return nil, err
	}
	kvs := make([]*config.KeyValue, 0, len(rsp.Kvs))
	for _, item := range rsp.Kvs {
		k := string(item.Key)
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
	return newWatcher(s), nil
}
