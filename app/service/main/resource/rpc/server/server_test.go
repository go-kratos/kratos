package server

import (
	"fmt"
	"net/rpc"
	"testing"
	"time"

	"go-common/app/service/main/resource/conf"
	"go-common/app/service/main/resource/model"
	"go-common/app/service/main/resource/service"
)

// rpc server const
const (
	addr           = "127.0.0.1:6429"
	_resourceAll   = "RPC.ResourceAll"
	_assignmentAll = "RPC.AssignmentAll"
	_defBanner     = "RPC.DefBanner"
	_resource      = "RPC.Resource"
	_resources     = "RPC.Resources"
	_assignment    = "RPC.Assignment"
	_banners       = "RPC.Banners"
	_pasterAPP     = "RPC.PasterAPP"
	_indexIcon     = "RPC.IndexIcon"
	_playerIcon    = "RPC.playerIcon"
	_cmtbox        = "RPC.Cmtbox"
	_sidebars      = "RPC.SideBars"
	_abtest        = "RPC.AbTest"
	_pasterCID     = "RPC.PasterCID"
)

// TestResource test rpc server
func TestResource(t *testing.T) {
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
	resourceAllRPC(client, t)
	assignmentAllRPC(client, t)
	defBannerRPC(client, t)
	resourceRPC(client, t)
	resourcesRPC(client, t)
	assignmentRPC(client, t)
	bannersRPC(client, t)
	pasterAPPRpc(client, t)
	indexIconRPC(client, t)
	playerIconRPC(client, t)
	cmtboxRPC(client, t)
	sideBarsRPC(client, t)
	abTestRPC(client, t)
	pasterCIDRPC(client, t)
}

func resourceAllRPC(client *rpc.Client, t *testing.T) {
	var res []*model.Resource
	arg := &struct{}{}
	if err := client.Call(_resourceAll, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("resourceAll", t, res)
	}
}

func assignmentAllRPC(client *rpc.Client, t *testing.T) {
	var res []*model.Assignment
	arg := &struct{}{}
	if err := client.Call(_assignmentAll, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("assignmentAll", t, res)
	}
}

func defBannerRPC(client *rpc.Client, t *testing.T) {
	var res model.Assignment
	arg := &struct{}{}
	if err := client.Call(_defBanner, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("defBanner", t, res)
	}
}

func resourceRPC(client *rpc.Client, t *testing.T) {
	var res model.Resource
	arg := &model.ArgRes{
		ResID: 1187,
	}
	if err := client.Call(_resource, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("resource", t, res)
	}
}

func resourcesRPC(client *rpc.Client, t *testing.T) {
	var res map[int]*model.Resource
	arg := &model.ArgRess{
		ResIDs: []int{1187, 1639},
	}
	if err := client.Call(_resources, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("resources", t, res)
	}
}

func assignmentRPC(client *rpc.Client, t *testing.T) {
	var res []*model.Assignment
	arg := &model.ArgRes{
		ResID: 1187,
	}
	if err := client.Call(_assignment, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("assignment", t, res)
	}
}

func bannersRPC(client *rpc.Client, t *testing.T) {
	var res *model.Banners
	arg := &model.ArgBanner{
		Plat:      1,
		ResIDs:    "454,467",
		Build:     508000,
		MID:       1493031,
		Channel:   "abc",
		IP:        "211.139.80.6",
		Buvid:     "123",
		Network:   "wifi",
		MobiApp:   "iphone",
		Device:    "test",
		IsAd:      true,
		OpenEvent: "abc",
	}
	if err := client.Call(_banners, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("banners", t, res)
	}
}

func pasterAPPRpc(client *rpc.Client, t *testing.T) {
	var res model.Paster
	arg := &model.ArgPaster{
		Platform: int8(1),
		AdType:   int8(1),
		Aid:      "666666",
		TypeId:   "11",
		Buvid:    "666666",
	}
	if err := client.Call(_pasterAPP, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("pasterAPPRpc", t, res)
	}
}

func indexIconRPC(client *rpc.Client, t *testing.T) {
	var res map[string][]*model.IndexIcon
	arg := &struct{}{}
	if err := client.Call(_indexIcon, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("indexIconRpc", t, res)
	}
}

func playerIconRPC(client *rpc.Client, t *testing.T) {
	var res *model.PlayerIcon
	arg := &struct{}{}
	if err := client.Call(_playerIcon, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("playerIconRPC", t, res)
	}
}

func cmtboxRPC(client *rpc.Client, t *testing.T) {
	var res model.Cmtbox
	arg := &model.ArgCmtbox{
		ID: 1,
	}
	if err := client.Call(_cmtbox, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("resource", t, res)
	}
}

func sideBarsRPC(client *rpc.Client, t *testing.T) {
	var res []*model.SideBars
	arg := &struct{}{}
	if err := client.Call(_sidebars, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("sideBars", t, res)
	}
}

func abTestRPC(client *rpc.Client, t *testing.T) {
	var res map[string]*model.AbTest
	arg := &model.ArgAbTest{
		Groups: "不显示热门tab,显示热门tab",
		IP:     "127.0.0.1",
	}
	if err := client.Call(_abtest, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("abTest", t, res)
	}
}

func pasterCIDRPC(client *rpc.Client, t *testing.T) {
	var res map[int64]int64
	arg := &struct{}{}
	if err := client.Call(_pasterCID, arg, &res); err != nil {
		t.Errorf("err: %v.", err)
	} else {
		result("pasterCID", t, res)
	}
}

func result(name string, t *testing.T, res interface{}) {
	fmt.Printf("res : %+v \n", res)
	t.Log("[==========" + name + "单元测试结果==========]")
	t.Log(res)
	t.Log("[↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑]\r\n")
}
