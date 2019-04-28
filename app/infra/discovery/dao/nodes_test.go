package dao

import (
	"context"
	"testing"
	"time"

	dc "go-common/app/infra/discovery/conf"
	"go-common/app/infra/discovery/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestReplicate(t *testing.T) {
	Convey("test replicate", t, func() {
		i := model.NewInstance(reg)
		nodes := NewNodes(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"api.bilibili.co", "uat-bilibili.co", "127.0.0.1:7171"}})
		nodes.nodes[0].client.SetTransport(gock.DefaultTransport)
		nodes.nodes[1].client.SetTransport(gock.DefaultTransport)
		httpMock("POST", "http://api.bilibili.co/discovery/register").Reply(200).JSON(`{"code":0}`)
		httpMock("POST", "http://uat-bilibili.co/discovery/register").Reply(200).JSON(`{"code":0}`)
		err := nodes.Replicate(context.TODO(), model.Register, i, false)
		So(err, ShouldBeNil)
	})
}

func TestReplicateSet(t *testing.T) {
	Convey("test replicate set", t, func() {
		nodes := NewNodes(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"api.bilibili.co"}})
		nodes.nodes[0].client.SetTransport(gock.DefaultTransport)
		httpMock("POST", "http://api.bilibili.co/discovery/set").Reply(200).JSON(`{"code":0}`)
		set := &model.ArgSet{
			Region:   "shsb",
			Env:      "pre",
			Appid:    "main.arch.account-service",
			Hostname: []string{"test1"},
			Status:   []int64{1},
		}
		err := nodes.ReplicateSet(context.TODO(), set, false)
		So(err, ShouldBeNil)
	})
}

func TestNodes(t *testing.T) {
	Convey("test replicate set", t, func() {
		nodes := NewNodes(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"api.bilibili.co", "uat-bilibili.co", "127.0.0.1:7171"}})
		res := nodes.Nodes()
		So(len(res), ShouldResemble, 3)
	})
}

func TestUp(t *testing.T) {
	Convey("test up", t, func() {
		nodes := NewNodes(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"api.bilibili.co", "uat-bilibili.co", "127.0.0.1:7171"}})
		nodes.UP()
		for _, nd := range nodes.nodes {
			if nd.addr == "127.0.0.1:7171" {
				So(nd.status, ShouldResemble, model.NodeStatusUP)
			}
		}
	})
}
