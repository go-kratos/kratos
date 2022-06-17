package config

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type Option func(*options)

type options struct {
	endpoint string // nolint:structcheck,unused

	namespaceID string // nolint:structcheck,unused

	group  string
	dataID string

	timeoutMs uint64
	logLevel  string

	logDir   string
	cacheDir string
}

// WithGroup With nacos config group.
func WithGroup(group string) Option {
	return func(o *options) {
		o.group = group
	}
}

// WithDataID With nacos config data id.
func WithDataID(dataID string) Option {
	return func(o *options) {
		o.dataID = dataID
	}
}

// WithLogDir With nacos config group.
func WithLogDir(logDir string) Option {
	return func(o *options) {
		o.logDir = logDir
	}
}

// WithCacheDir With nacos config cache dir.
func WithCacheDir(cacheDir string) Option {
	return func(o *options) {
		o.cacheDir = cacheDir
	}
}

// WithLogLevel With nacos config log level.
func WithLogLevel(logLevel string) Option {
	return func(o *options) {
		o.logLevel = logLevel
	}
}

// WithTimeout With nacos config timeout.
func WithTimeout(time time.Duration) Option {
	return func(o *options) {
		o.timeoutMs = uint64(time.Milliseconds())
	}
}

type Config struct {
	opts   options
	client config_client.IConfigClient
}

func NewConfigSource(client config_client.IConfigClient, opts ...Option) config.Source {
	_options := options{}
	for _, o := range opts {
		o(&_options)
	}
	return &Config{client: client, opts: _options}
}

func (c *Config) Load() ([]*config.KeyValue, error) {
	content, err := c.client.GetConfig(vo.ConfigParam{
		DataId: c.opts.dataID,
		Group:  c.opts.group,
	})
	if err != nil {
		return nil, err
	}
	k := c.opts.dataID
	return []*config.KeyValue{
		{
			Key:    k,
			Value:  []byte(content),
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		},
	}, nil
}

func (c *Config) Watch() (config.Watcher, error) {
	watcher := newWatcher(context.Background(), c.opts.dataID, c.opts.group, c.client.CancelListenConfig)
	err := c.client.ListenConfig(vo.ConfigParam{
		DataId: c.opts.dataID,
		Group:  c.opts.group,
		OnChange: func(_, group, dataId, data string) {
			if dataId == watcher.dataID && group == watcher.group {
				watcher.content <- data
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return watcher, nil
}
