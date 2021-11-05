package config

import (
	"fmt"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestGetConfig(t *testing.T) {
	ip := ""
	//ctx := context.Background()

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "5c7d1f8b-6782-46bf-b36a-dfd9cfc14b89", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// a more graceful way to create naming client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	c := kconfig.New(
		kconfig.WithSource(
			NewConfigSource(client, WithGroup("private"), WithDataID("go-rpc-executor")),
		),
		kconfig.WithDecoder(func(kv *kconfig.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	name, err := c.Value("logger.level").String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("get value",name)
	// 监听值内容变更
	done := make(chan bool)
	c.Watch("logger.level", func(key string, value kconfig.Value) {
		fmt.Println(key ," value change", value)
		done <- true
	})
	<-done
}