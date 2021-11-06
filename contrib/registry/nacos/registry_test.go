package nacos

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/go-kratos/kratos/v2/registry"
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

func TestRegistry(t *testing.T) {
	ip := getIntranetIP()
	serviceName := "golang-sms@grpc"
	ctx := context.Background()

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public", // namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	is, err := client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("is %#v", is)

	time.Sleep(time.Second)
	r := New(client)

	go func() {
		var (
			w        registry.Watcher
			watchErr error
		)
		w, watchErr = r.Watch(ctx, "golang-sms@grpc")
		if watchErr != nil {
			log.Fatal(watchErr)
		}
		for {
			var res []*registry.ServiceInstance
			res, watchErr = w.Next()
			if watchErr != nil {
				return
			}
			log.Printf("watch: %d", len(res))
			for _, r := range res {
				log.Printf("next: %+v", r)
			}
		}
	}()

	time.Sleep(time.Second)

	ins, err := r.GetService(ctx, serviceName)
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range ins {
		t.Logf("ins: %#v", in)
	}

	time.Sleep(time.Second)
}

func TestRegistryMany(t *testing.T) {
	ip := getIntranetIP()
	serviceName := "golang-sms@grpc"
	// ctx := context.Background()

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public", // namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f1",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f2",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "f3",
		Port:        8840,
		ServiceName: serviceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai-xs"},
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	is, err := client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, host := range is.Hosts {
		t.Logf("host: %#v,e: %v", host, err)
	}

	time.Sleep(time.Second)
}
