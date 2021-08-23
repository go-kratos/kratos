package registry

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func tcpServer(t *testing.T, lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		fmt.Println("get tcp")
		conn.Close()
	}
}

func TestRegister(t *testing.T) {
	addr := fmt.Sprintf("%s:8081", getIntranetIP())
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("listen tcp %s failed!", addr)
		t.Fail()
	}
	defer lis.Close()
	go tcpServer(t, lis)
	time.Sleep(time.Millisecond * 100)
	cli, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}
	r := New(cli)
	assert.Nil(t, err)
	version := strconv.FormatInt(time.Now().Unix(), 10)
	svc := &registry.ServiceInstance{
		ID:        "test2233",
		Name:      "test-provider",
		Version:   version,
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = r.Deregister(ctx, svc)
	assert.Nil(t, err)
	err = r.Register(ctx, svc)
	assert.Nil(t, err)
	w, err := r.Watch(ctx, "test-provider")
	assert.Nil(t, err)

	services, err := w.Next()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.EqualValues(t, "test2233", services[0].ID)
	assert.EqualValues(t, "test-provider", services[0].Name)
	assert.EqualValues(t, version, services[0].Version)
}

func getIntranetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
