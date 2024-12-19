package nacos

import (
	"context"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	REGISTRY_GROUP   = "DEFAULT_GROUP"
	REGISTRY_CLUSTER = "DEFAULT"
)

func newNamingClient() (naming_client.INamingClient, error) {
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

	return clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
}

func TestRegistry(t *testing.T) {
	client, err := newNamingClient()
	if err != nil {
		t.Fatal(err)
	}

	nacosRegistry := NewRegistry(
		client,
		WithRegistryGroup(REGISTRY_GROUP),
		WithCluster(REGISTRY_CLUSTER),
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

	nacosRegistry.Register(context.Background(), ins)

	instances, err := nacosWatcher.Next()
	for _, instance := range instances {
		t.Log(instance.Version)
	}
}
