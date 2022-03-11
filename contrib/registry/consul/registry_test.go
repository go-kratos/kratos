package consul

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/SeeMusic/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
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
	addr := fmt.Sprintf("%s:9091", getIntranetIP())
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
	opts := []Option{
		WithHeartbeat(true),
		WithHealthCheck(true),
		WithHealthCheckInterval(5),
	}
	r := New(cli, opts...)
	if err != nil {
		t.Errorf("new consul registry failed: %v", err)
	}
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
	if err != nil {
		t.Errorf("Deregister failed: %v", err)
	}
	err = r.Register(ctx, svc)
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}
	w, err := r.Watch(ctx, "test-provider")
	if err != nil {
		t.Errorf("Watchfailed: %v", err)
	}

	services, err := w.Next()
	if err != nil {
		t.Errorf("Next failed: %v", err)
	}
	if !reflect.DeepEqual(1, len(services)) {
		t.Errorf("no expect float_key value: %v, but got: %v", len(services), 1)
	}
	if !reflect.DeepEqual("test2233", services[0].ID) {
		t.Errorf("no expect float_key value: %v, but got: %v", services[0].ID, "test2233")
	}
	if !reflect.DeepEqual("test-provider", services[0].Name) {
		t.Errorf("no expect float_key value: %v, but got: %v", services[0].Name, "test-provider")
	}
	if !reflect.DeepEqual(version, services[0].Version) {
		t.Errorf("no expect float_key value: %v, but got: %v", services[0].Version, version)
	}
}

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
