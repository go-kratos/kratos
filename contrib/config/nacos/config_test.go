package config

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v3"

	kconfig "github.com/SeeMusic/kratos/v2/config"
)

func getIntranetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func TestGetConfig(t *testing.T) {
	ip := getIntranetIP()
	// ctx := context.Background()

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	cc := constant.ClientConfig{
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

	dataID := "test.yaml"
	group := "test"
	_, err = client.PublishConfig(vo.ConfigParam{DataId: dataID, Group: group, Content: `
logger:
  level: info
`})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	c := kconfig.New(
		kconfig.WithSource(
			NewConfigSource(client, WithGroup(group), WithDataID(dataID)),
		),
		kconfig.WithDecoder(func(kv *kconfig.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)

	if err = c.Load(); err != nil {
		t.Fatal(err)
	}

	name, err := c.Value("logger.level").String()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("get value", name)

	done := make(chan struct{})
	err = c.Watch("logger.level", func(key string, value kconfig.Value) {
		fmt.Println(key, " value change", value)
		done <- struct{}{}
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.PublishConfig(vo.ConfigParam{DataId: dataID, Group: group, Content: `
logger:
  level: debug
`})
	if err != nil {
		t.Fatal(err)
	}

	<-done
}
