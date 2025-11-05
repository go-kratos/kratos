package nacos

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	DataID = "nacos-test-config.json"
	Group  = "DEFAULT_GROUP"
)

func newConfigClient() (config_client.IConfigClient, error) {
	userHomeDir, _ := os.UserHomeDir()

	// 设置默认值
	serverAddr := os.Getenv("NACOS_SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "127.0.0.1" // 设置默认地址
	}

	username := os.Getenv("NACOS_USER_NAME")
	password := os.Getenv("NACOS_USER_PASSWORD")

	clientConfig := constant.NewClientConfig(
		constant.WithLogDir(filepath.Join(userHomeDir, "logs", "nacos")),
		constant.WithCacheDir(filepath.Join(userHomeDir, "nacos", "cache")),
		constant.WithUsername(username),
		constant.WithPassword(password),
		constant.WithLogLevel("info"),
	)

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      serverAddr,
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http", // 添加scheme
		},
	}

	return clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
}

func TestConfig_Get_AND_WATCH(t *testing.T) {
	// 检查必要的环境变量
	if os.Getenv("NACOS_SERVER_ADDRESS") == "" {
		t.Skip("NACOS_SERVER_ADDRESS environment variable not set")
	}

	client, err := newConfigClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  DataID,
		Group:   Group,
		Content: "hello world",
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	source := NewConfigSource(client, WithConfigGroup(Group), WithDataID(DataID))

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
		DataId:  DataID,
		Group:   Group,
		Content: "hello world2",
	})
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
