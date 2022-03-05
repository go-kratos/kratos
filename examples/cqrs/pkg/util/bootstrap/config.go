package bootstrap

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	// etcd config
	etcdKratos "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	etcdV3 "go.etcd.io/etcd/client/v3"
	GRPC "google.golang.org/grpc"

	// consul config
	consulKratos "github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/hashicorp/consul/api"

	// nacos config
	nacosKratos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	nacosClients "github.com/nacos-group/nacos-sdk-go/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// getConfigKey 获取合法的配置名
func getConfigKey(configKey string, useBackslash bool) string {
	if useBackslash {
		return strings.Replace(configKey, `.`, `/`, -1)
	} else {
		return configKey
	}
}

// NewRemoteConfigSource 创建一个远程配置源
func NewRemoteConfigSource(configType, configHost, configKey string) config.Source {
	switch configType {
	case "nacos":
		uri, _ := url.Parse(configHost)
		h := strings.Split(uri.Host, ":")
		addr := h[0]
		port, _ := strconv.Atoi(h[1])
		return NewNacosConfigSource(addr, uint64(port), configKey)
	case "consul":
		return NewConsulConfigSource(configHost, configKey)
	case "etcd":
		return NewEtcdConfigSource(configHost, configKey)
	case "apollo":
		return NewApolloConfigSource(configHost, configKey)
	}
	return nil
}

// NewNacosConfigSource 创建一个远程配置源 - Nacos
func NewNacosConfigSource(configAddr string, configPort uint64, configKey string) config.Source {
	sc := []nacosConstant.ServerConfig{
		*nacosConstant.NewServerConfig(configAddr, configPort),
	}

	cc := nacosConstant.ClientConfig{
		TimeoutMs:            10 * 1000, // http请求超时时间，单位毫秒
		BeatInterval:         5 * 1000,  // 心跳间隔时间，单位毫秒
		UpdateThreadNum:      20,        // 更新服务的线程数
		LogLevel:             "debug",
		CacheDir:             "../../configs/cache", // 缓存目录
		LogDir:               "../../configs/log",   // 日志目录
		NotLoadCacheAtStart:  true,                  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: true,                  // 当服务列表为空时是否更新本地缓存，true--更新,false--不更新
	}

	nacosClient, err := nacosClients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	return nacosKratos.NewConfigSource(nacosClient,
		nacosKratos.WithGroup(getConfigKey(configKey, false)),
		nacosKratos.WithDataID("bootstrap.yaml"),
	)
}

// NewEtcdConfigSource 创建一个远程配置源 - Etcd
func NewEtcdConfigSource(configHost, configKey string) config.Source {
	etcdClient, err := etcdV3.New(etcdV3.Config{
		Endpoints:   []string{configHost},
		DialTimeout: time.Second, DialOptions: []GRPC.DialOption{GRPC.WithBlock()},
	})
	if err != nil {
		panic(err)
	}

	etcdSource, err := etcdKratos.New(etcdClient, etcdKratos.WithPath(getConfigKey(configKey, true)))
	if err != nil {
		panic(err)
	}

	return etcdSource
}

// NewApolloConfigSource 创建一个远程配置源 - Apollo
func NewApolloConfigSource(_, _ string) config.Source {
	return nil
}

// NewConsulConfigSource 创建一个远程配置源 - Consul
func NewConsulConfigSource(configHost, configKey string) config.Source {
	consulClient, err := api.NewClient(&api.Config{
		Address: configHost,
	})
	if err != nil {
		panic(err)
	}

	consulSource, err := consulKratos.New(consulClient,
		consulKratos.WithPath(getConfigKey(configKey, true)),
	)
	if err != nil {
		panic(err)
	}

	return consulSource
}

// NewFileConfigSource 创建一个本地文件配置源
func NewFileConfigSource(filePath string) config.Source {
	return file.NewSource(filePath)
}

// NewConfigProvider 创建一个配置
func NewConfigProvider(configType, configHost, configPath, configKey string) config.Config {
	return config.New(
		config.WithSource(
			NewFileConfigSource(configPath),
			NewRemoteConfigSource(configType, configHost, configKey),
		),
	)
}
