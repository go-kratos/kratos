package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "go-common/app/service/main/location/api"
	"go-common/app/service/main/location/conf"
	"go-common/app/service/main/location/service"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

// rpc server const
const (
	addr = "127.0.0.1:9000"
)

// TestLocation test rpc server
func TestLocation(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf.WardenServer, svr)
	time.Sleep(time.Second * 3)
	cfg := &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 3),
		Timeout: xtime.Duration(time.Second * 3),
	}
	cc, err := warden.NewClient(cfg).Dial(context.Background(), addr)
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	client := pb.NewLocationClient(cc)
	infoGRPC(client, t)
}

func infoGRPC(client pb.LocationClient, t *testing.T) {
	arg := &pb.InfoReq{
		Addr: "211.139.80.6",
	}
	res, err := client.Info(context.TODO(), arg)
	if err != nil {
		t.Error(err)
	} else {
		result("info", t, res)
	}
}

func result(name string, t *testing.T, res interface{}) {
	fmt.Printf("res : %+v \n", res)
	t.Log("[==========" + name + "单元测试结果==========]")
	t.Log(res)
	t.Log("[↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑]\r\n")
}
