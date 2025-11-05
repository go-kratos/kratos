package nacos

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type ConfigOption func(*configOptions)

type configOptions struct {
	group    string
	dataID   string
	kmsKeyID string
}

func WithConfigGroup(group string) ConfigOption {
	return func(o *configOptions) {
		o.group = group
	}
}

func WithDataID(dataID string) ConfigOption {
	return func(o *configOptions) {
		o.dataID = dataID
	}
}

func WithKmsKeyID(kmsKeyID string) ConfigOption {
	return func(o *configOptions) {
		o.kmsKeyID = kmsKeyID
	}
}

type Config struct {
	opts   configOptions
	client config_client.IConfigClient
}

func NewConfigSource(client config_client.IConfigClient, opts ...ConfigOption) config.Source {
	_options := configOptions{}
	for _, o := range opts {
		o(&_options)
	}
	return &Config{client: client, opts: _options}
}

func (c *Config) Load() ([]*config.KeyValue, error) {
	content, err := c.client.GetConfig(vo.ConfigParam{
		DataId:   c.opts.dataID,
		Group:    c.opts.group,
		KmsKeyId: c.opts.kmsKeyID,
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
		DataId:   c.opts.dataID,
		Group:    c.opts.group,
		KmsKeyId: c.opts.kmsKeyID,
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
