package gorpc

import (
	"flag"
	"net/rpc"
	"testing"
	"time"

	"go-common/app/service/main/member/conf"
	"go-common/app/service/main/member/model"
	"go-common/app/service/main/member/service"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	flag.Set("conf", "../../cmd/member-service-example.toml")
	startService()
}

const (
	addr      = "127.0.0.1:6689"
	_testPing = "RPC.Ping"
)

var (
	_noArg = &struct{}{}
	svr    *service.Service
	client *rpc.Client
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	svr = service.New(conf.Conf)
	New(conf.Conf, svr)
	time.Sleep(time.Second * 3)
	var err error
	client, err = rpc.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
}
func TestAccountRpc(t *testing.T) {
	Convey("ping", t, func() {
		err := client.Call(_testPing, &_noArg, &_noArg)
		So(err, ShouldBeNil)
	})
}

func TestExp(t *testing.T) {
	Convey("update", t, func() {
		err := client.Call("RPC.UpdateExp", &model.ArgAddExp{
			Mid:     1,
			Count:   2,
			Reason:  "test",
			Operate: "other",
			IP:      "111",
		}, &_noArg)
		So(err, ShouldBeNil)

	})
	Convey("exp", t, func() {
		res := new(model.LevelInfo)
		err := client.Call("RPC.Exp", &model.ArgMid{Mid: 1}, res)
		So(err, ShouldBeNil)
		So(res.NextExp, ShouldNotEqual, 0)
	})
}
func TestLevel(t *testing.T) {
	Convey("level", t, func() {
		res := new(model.LevelInfo)
		err := client.Call("RPC.Level", &model.ArgMid{
			Mid: 1,
		}, res)
		So(err, ShouldNotBeNil)
		So(res.NextExp, ShouldNotEqual, 0)
	})

}

func TestLog(t *testing.T) {
	Convey("log", t, func() {
		var res []*model.UserLog
		err := client.Call("RPC.Log", &model.ArgMid{Mid: 1}, &res)
		So(err, ShouldNotBeNil)
	})
}
