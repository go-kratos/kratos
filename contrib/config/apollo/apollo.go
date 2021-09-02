package apollo

import (
	"github.com/go-kratos/kratos/v2/config"

	"github.com/apolloconfig/agollo/v4"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
)

type apollo struct {
	client *agollo.Client
}

type Option func(*options)

type options struct {
	ip             string
	appid          string
	cluster        string
	secret         string
	namespaceName  string
	isBackupConfig bool
}

// AppID with apollo config app id
func AppID(appID string) Option {
	return func(o *options) {
		o.appid = appID
	}
}

// Cluster with apollo config cluster
func Cluster(cluster string) Option {
	return func(o *options) {
		o.cluster = cluster
	}
}

// IP with apollo config conf server ip
func IP(ip string) Option {
	return func(o *options) {
		o.ip = ip
	}
}

// IsBackupConfig with apollo config isBackupConfig
func IsBackupConfig(isBackupConfig bool) Option {
	return func(o *options) {
		o.isBackupConfig = isBackupConfig
	}
}

// Secret with apollo config app secret
func Secret(secret string) Option {
	return func(o *options) {
		o.secret = secret
	}
}

// NamespaceName with apollo config namespace name
func NamespaceName(name string) Option {
	return func(o *options) {
		o.namespaceName = name
	}
}

// NewSource start with config file in ENV
// Linux/Mac export AGOLLO_CONF=/a/conf.properties
// Windows set AGOLLO_CONF=c:/a/conf.properties
// more detail:https://github.com/apolloconfig/agollo/wiki/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97#1312%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F%E6%8C%87%E5%AE%9A%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6
func NewSource() config.Source {
	client, err := agollo.Start()
	if err != nil {
		panic(err)
	}
	return &apollo{client}
}

func NewSourceWithConfig(opts ...Option) config.Source {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	client, err := agollo.StartWithConfig(func() (*apolloConfig.AppConfig, error) {
		return &apolloConfig.AppConfig{
			AppID:          op.appid,
			Cluster:        op.cluster,
			NamespaceName:  op.namespaceName,
			IP:             op.ip,
			IsBackupConfig: op.isBackupConfig,
			Secret:         op.secret,
		}, nil
	})
	if err != nil {
		panic(err)
	}
	return &apollo{client}
}

func (e *apollo) load() []*config.KeyValue {
	kv := make([]*config.KeyValue, 0)
	e.client.GetDefaultConfigCache().Range(func(key, value interface{}) bool {
		kv = append(kv, &config.KeyValue{
			Key:   key.(string),
			Value: []byte(value.(string)),
		})
		return true
	})
	return kv
}

func (e *apollo) Load() (kv []*config.KeyValue, err error) {
	return e.load(), nil
}

func (e *apollo) Watch() (config.Watcher, error) {
	w, err := NewWatcher(e)
	if err != nil {
		return nil, err
	}
	return w, nil
}
