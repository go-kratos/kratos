package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	DATA_ID = "nacos-test-config.json"
	GROUP   = "DEFAULT_GROUP"
)

func newConfigClient() (config_client.IConfigClient, error) {
	userHomeDir, _ := os.UserHomeDir()

	clientConfig := constant.NewClientConfig(
		constant.WithLogDir(filepath.Join(userHomeDir, "logs", "nacos")),
		constant.WithCacheDir(filepath.Join(userHomeDir, "nacos", "cache")),
		constant.WithUsername(os.Getenv("NACOS_USER_NAME")),
		constant.WithPassword(os.Getenv("NACOS_USER_PASSWORD")),
		constant.WithLogLevel("info"),
	)
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      os.Getenv("NACOS_SERVER_ADDRESS"),
			ContextPath: "/nacos",
			Port:        8848,
		},
	}

	return clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
}

func TestConfig_Get_AND_WATCH(t *testing.T) {
	client, err := newConfigClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  DATA_ID,
		Group:   GROUP,
		Content: "hello world"})

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	source := NewConfigSource(client, WithConfigGroup(GROUP), WithDataID(DATA_ID))

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, kv := range kvs {
		t.Logf("key: %s, value: %s", kv.Key, kv.Value)
	}

	w, err := source.Watch()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  DATA_ID,
		Group:   GROUP,
		Content: "hello world2"})

	if err != nil {
		t.Fatal(err)
	}

	next, err := w.Next()
	if err != nil {
		t.Fatal(err)
	}

	for _, kv := range next {
		t.Logf("key: %s, value: %s", kv.Key, kv.Value)
	}
}
