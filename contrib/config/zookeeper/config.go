package zookeeper

import (
	"context"
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-zookeeper/zk"
)

// Option is zookeeper config option.
type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	key       string
	user      string
	password  string
}

// WithContext with config context.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithNamespace is config path namespace
func WithNamespace(n string) Option {
	return func(o *options) {
		o.namespace = n
	}
}

// WithKey is config path key
func WithKey(k string) Option {
	return func(o *options) {
		o.key = k
	}
}

// WithDigestACL with zookeeper username and password.
func WithDigestACL(user string, password string) Option {
	return func(o *options) {
		o.user = user
		o.password = password
	}
}

type source struct {
	conn    *zk.Conn
	options *options
}

func New(conn *zk.Conn, opts ...Option) (config.Source, error) {
	options := &options{
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.namespace == "" {
		return nil, errors.New("path invalid")
	}

	if options.key == "" {
		return nil, errors.New("key invalid")
	}

	return &source{
		conn:    conn,
		options: options,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	fullPath := path.Join(s.options.namespace, s.options.key)
	res, _, err := s.conn.Get(fullPath)
	if err != nil {
		return nil, err
	}

	return []*config.KeyValue{{
		Key:    s.options.key,
		Value:  res,
		Format: strings.TrimPrefix(filepath.Ext(s.options.key), "."),
	}}, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s), nil
}
