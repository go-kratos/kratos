package nacos

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/naming"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	//ctx := context.TODO()
	//s1 := createServer("server1", "127.0.0.1:18001")
	//s2 := createServer("server2", "127.0.0.1:18002")
	//defer s1.Shutdown(ctx)
	//defer s2.Shutdown(ctx)
	os.Exit(m.Run())
}

func TestNacos(t *testing.T) {
	config := &Config{
		ServerConfigs: []constant.ServerConfig{
			{
				IpAddr: "192.168.9.102",
				Port:   8848,
			},
		},
		ClientConfig: constant.ClientConfig{
			TimeoutMs:           10 * 1000,
			BeatInterval:        5 * 1000,
			ListenInterval:      30 * 1000,
			NotLoadCacheAtStart: true},
	}
	nacos, err := New(config)
	if err != nil {
		panic(err)
	}
	//
	instance := &naming.Instance{
		Region:   "china",
		Zone:     "shanghai",
		Env:      "dev",
		AppID:    "test-nacos",
		Hostname: "",
		Addrs:    []string{"grpc://127.0.0.1:18080", "grpc://127.0.0.2:18080"},
		Version:  "v3.0",
		Metadata: map[string]string{"weight": "10"},
		Status:   1,
	}
	//
	go func() {
		resolver := nacos.Build("test-nacos")
		for true {
			instancesInfo, success := resolver.Fetch(context.Background())
			if !success {
				t.Error("fetch failed")
			} else {
				for _, instances := range instancesInfo.Instances {
					for _, i := range instances {
						fmt.Println(i)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	cancelFunc, err := nacos.Register(context.Background(), instance)
	if err != nil {
		panic(err)
	}
	var c chan int
	c <- 1
	cancelFunc()
}
