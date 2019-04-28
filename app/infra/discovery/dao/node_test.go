package dao

import (
	"context"
	"strings"
	"testing"
	"time"

	dc "go-common/app/infra/discovery/conf"
	"go-common/app/infra/discovery/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestCall(t *testing.T) {
	Convey("test call", t, func() {
		var res *model.Instance
		node := newNode(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"127.0.0.1:7171"}}, "api.bilibili.co")
		node.client.SetTransport(gock.DefaultTransport)
		httpMock("POST", "http://api.bilibili.co/discovery/register").Reply(200).JSON(`{"ts":1514341945,"code":-409,"data":{"region":"shsb","zone":"fuck","appid":"main.arch.account-service","env":"pre","hostname":"cs4sq","http":"","rpc":"0.0.0.0:18888","weight":2}}`)
		i := model.NewInstance(reg)
		err := node.call(context.TODO(), model.Register, i, "http://api.bilibili.co/discovery/register", &res)
		So(err, ShouldResemble, ecode.Conflict)
		So(res.Appid, ShouldResemble, "main.arch.account-service")
	})
}

func TestNodeCancel(t *testing.T) {
	Convey("test node renew 409 error", t, func() {
		i := model.NewInstance(reg)
		node := newNode(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"127.0.0.1:7171"}}, "api.bilibili.co")
		node.pRegisterURL = "http://127.0.0.1:7171/discovery/register"
		node.client.SetTransport(gock.DefaultTransport)
		httpMock("POST", "http://api.bilibili.co/discovery/cancel").Reply(200).JSON(`{"code":0}`)
		err := node.Cancel(context.TODO(), i)
		So(err, ShouldBeNil)
	})
}

func TestNodeRenew(t *testing.T) {
	Convey("test node renew 409 error", t, func() {
		i := model.NewInstance(reg)
		node := newNode(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"127.0.0.1:7171"}}, "api.bilibili.co")
		node.pRegisterURL = "http://127.0.0.1:7171/discovery/register"
		node.client.SetTransport(gock.DefaultTransport)
		httpMock("POST", "http://api.bilibili.co/discovery/renew").Reply(200).JSON(`{"code":-409,"data":{"region":"shsb","zone":"fuck","appid":"main.arch.account-service","env":"pre","hostname":"cs4sq","http":"","rpc":"0.0.0.0:18888","weight":2}}`)
		httpMock("POST", "http://127.0.0.1:7171/discovery/register").Reply(200).JSON(`{"code":0}`)
		err := node.Renew(context.TODO(), i)
		So(err, ShouldBeNil)
	})
}

func TestNodeRenew2(t *testing.T) {
	Convey("test node renew 404 error", t, func() {
		i := model.NewInstance(reg)
		node := newNode(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"127.0.0.1:7171"}}, "api.bilibili.co")
		node.client.SetTransport(gock.DefaultTransport)
		httpMock("POST", "http://api.bilibili.co/discovery/renew").Reply(200).JSON(`{"code":-404}`)
		httpMock("POST", "http://api.bilibili.co/discovery/register").Reply(200).JSON(`{"code":0}`)
		err := node.Renew(context.TODO(), i)
		So(err, ShouldBeNil)
	})
}

func TestSet(t *testing.T) {
	Convey("test set", t, func() {
		node := newNode(&dc.Config{HTTPClient: &bm.ClientConfig{Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100}, Timeout: xtime.Duration(time.Second), App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}}, BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}}, Nodes: []string{"127.0.0.1:7171"}}, "api.bilibili.co")
		node.client.SetTransport(gock.DefaultTransport)
		httpMock("POST", "http://api.bilibili.co/discovery/set").Reply(200).JSON(`{"ts":1514341945,"code":0}`)
		set := &model.ArgSet{
			Region:   "shsb",
			Env:      "pre",
			Appid:    "main.arch.account-service",
			Hostname: []string{"test1"},
			Status:   []int64{1},
		}
		err := node.Set(context.TODO(), set)
		So(err, ShouldBeNil)
	})
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
