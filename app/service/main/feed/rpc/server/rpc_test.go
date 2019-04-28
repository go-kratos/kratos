package server

import (
	artmdl "go-common/app/interface/openplatform/article/model"
	feed "go-common/app/service/main/feed/model"
	"net/rpc"
	"testing"
)

const (
	addr = "172.16.33.57:6361"

	_testArticleFeed = "RPC.ArticleFeed"
)

func TestFeedRpc(t *testing.T) {
	client, err := rpc.Dial("tcp", addr)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial(tcp, \"%s\") error(%v)", addr, err)
		t.FailNow()
	}
	feedRPC(client, t)
}

func feedRPC(client *rpc.Client, t *testing.T) {
	arg := &feed.ArgFeed{}
	arg.Mid = 88888929
	res := &[]*artmdl.Meta{}
	if err := client.Call(_testArticleFeed, arg, &res); err != nil {
		t.Errorf("client.Call(%s) error(%v)", _testArticleFeed, err)
		t.FailNow()
	} else {
		result("article", t, res)
	}
}

func result(name string, t *testing.T, res interface{}) {
	t.Log("[==========" + name + "单元测试结果==========]")
	t.Log(res)
	t.Log("[↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑]\r\n")
}
