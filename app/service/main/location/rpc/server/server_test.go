package server

import (
	"fmt"
	"net/rpc"
	"testing"
	"time"

	"go-common/app/service/main/location/conf"
	"go-common/app/service/main/location/model"
	"go-common/app/service/main/location/service"
)

// rpc server const
const (
	addr           = "127.0.0.1:6293"
	_archive       = "RPC.Archive"
	_archive2      = "RPC.Archive2"
	_group         = "RPC.Group"
	_authPIDs      = "RPC.AuthPIDs"
	_info          = "RPC.Info"
	_infos         = "RPC.Infos"
	_infoComplete  = "RPC.InfoComplete"
	_infosComplete = "RPC.InfosComplete"
)

// TestLocation test rpc server
func TestLocation(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)
	time.Sleep(time.Second * 3)

	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	archiveRPC(client, t)
	archive2RPC(client, t)
	groupRPC(client, t)
	authPIDsRPC(client, t)
	infoRPC(client, t)
	infosRPC(client, t)
	infoCompleteRPC(client, t)
	infosCompleteRPC(client, t)
}

func archiveRPC(client *rpc.Client, t *testing.T) {
	var res int64
	arg := &model.Archive{
		Aid: 740955,
		Mid: 0,
		IP:  "2.20.32.123",
		CIP: "127.0.0.1",
	}
	err := client.Call(_archive, arg, &res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("archive", t, res)
	}
	t.Logf("%+v", res)
}

func archive2RPC(client *rpc.Client, t *testing.T) {
	res := &model.Auth{}
	arg := &model.Archive{
		Aid: 740955,
		Mid: 0,
		IP:  "2.20.32.123",
		CIP: "127.0.0.1",
	}
	err := client.Call(_archive2, arg, &res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("archive2", t, res)
	}
	t.Logf("%+v", res)
}

func groupRPC(client *rpc.Client, t *testing.T) {
	res := &model.Auth{}
	arg := &model.Group{
		Gid: 317,
		Mid: 0,
		IP:  "2.20.32.123",
		CIP: "127.0.0.1",
	}
	err := client.Call(_group, arg, &res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("group", t, res)
	}
	t.Logf("%+v", res)
}

func authPIDsRPC(client *rpc.Client, t *testing.T) {
	arg := &model.ArgPids{IP: "61.216.166.156", Pids: "150,92"}
	res := map[int64]*model.Auth{}
	err := client.Call(_authPIDs, arg, &res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("authPIDs", t, res)
	}
	t.Logf("%+v,err:%v.", res, err)
}

func infoRPC(client *rpc.Client, t *testing.T) {
	var res = new(model.Info)
	arg := &model.ArgIP{
		IP: "139.214.144.59",
	}
	err := client.Call(_info, arg, res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("info", t, res)
	}
	t.Logf("%+v,err:%v.", res, err)
}

func infosRPC(client *rpc.Client, t *testing.T) {
	arg := []string{"61.216.166.156", "211.139.80.6"}
	res := map[string]*model.Info{}
	err := client.Call(_infos, arg, &res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("infos", t, res)
	}
	t.Logf("%+v,err:%v.", res, err)
}

func infoCompleteRPC(client *rpc.Client, t *testing.T) {
	var res = new(model.InfoComplete)
	arg := &model.ArgIP{
		IP: "139.214.144.59",
	}
	err := client.Call(_infoComplete, arg, res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("infoComplete", t, res)
	}
	t.Logf("%+v,err:%v.", res, err)
}

func infosCompleteRPC(client *rpc.Client, t *testing.T) {
	arg := []string{"61.216.166.156", "211.139.80.6"}
	res := map[string]*model.InfoComplete{}
	err := client.Call(_infosComplete, arg, &res)
	if err != nil {
		t.Errorf("err:%v.", err)
	} else {
		result("infosComplete", t, res)
	}
	t.Logf("%+v,err:%v.", res, err)
}

func result(name string, t *testing.T, res interface{}) {
	fmt.Printf("res : %+v \n", res)
	t.Log("[==========" + name + "单元测试结果==========]")
	t.Log(res)
	t.Log("[↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑]\r\n")
}
