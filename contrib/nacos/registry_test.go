package nacos

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/go-kratos/kratos/v2/registry"
)

var (
	RegistryGroup   = "DEFAULT_GROUP"
	RegistryCluster = "DEFAULT"
)

func newNamingClient() (naming_client.INamingClient, error) {
	userHomeDir, _ := os.UserHomeDir()

	serverAddr := os.Getenv("NACOS_SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "127.0.0.1"
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
			Scheme:      "http",
		},
	}

	return clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
}

func TestRegistry(t *testing.T) {
	if os.Getenv("NACOS_SERVER_ADDRESS") == "" {
		t.Skip("NACOS_SERVER_ADDRESS environment variable not set")
	}

	client, err := newNamingClient()
	if err != nil {
		t.Fatal(err)
	}

	nacosRegistry := NewRegistry(
		client,
		WithRegistryGroup(RegistryGroup),
		WithCluster(RegistryCluster),
		WithWeight(1.0))

	mm := map[string]string{
		"test1": "test1",
	}
	ins := &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
		},
		Metadata: mm,
	}

	err = nacosRegistry.Register(context.Background(), ins)
	if err != nil {
		t.Fatal(err)
	}
	nacosWatcher, err := nacosRegistry.Watch(context.Background(), "test-ut")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err = nacosRegistry.Deregister(context.Background(), ins); err != nil {
			t.Fatal(err)
		}

		if err = nacosWatcher.Stop(); err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	service, err := nacosRegistry.GetService(context.Background(), "test-ut")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)

	ins = &registry.ServiceInstance{
		ID:      "test-ut",
		Name:    "test-ut",
		Version: "v2.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
		},
		Metadata: mm,
	}

	err = nacosRegistry.Register(context.Background(), ins)
	if err != nil {
		return
	}

	instances, err := nacosWatcher.Next()
	for _, instance := range instances {
		t.Log(instance.Version)
	}
}
