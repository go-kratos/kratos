package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	dc "go-common/app/infra/discovery/conf"
	"go-common/app/infra/discovery/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx     = context.TODO()
	reg     = defRegisArg()
	rew     = &model.ArgRenew{Appid: "main.arch.test", Hostname: "test1", Region: "shsb", Zone: "sh001", Env: "pre"}
	cancel  = &model.ArgCancel{Appid: "main.arch.test", Hostname: "test1", Region: "shsb", Zone: "sh001", Env: "pre"}
	fet     = &model.ArgFetch{Appid: "main.arch.test", Region: "shsb", Zone: "sh001", Env: "pre", Status: 1}
	set     = &model.ArgSet{Appid: "main.arch.test", Region: "shsb", Hostname: []string{"test1"}, Zone: "sh001", Env: "pre"}
	pollArg = newPoll()
)

func newFetchArg() *model.ArgFetchs {
	return &model.ArgFetchs{Appid: []string{"main.arch.test"}, Zone: "sh001", Env: "pre", Status: 1}
}
func newPoll() *model.ArgPolls {
	return &model.ArgPolls{
		Region:          "shsb",
		Env:             "pre",
		Appid:           []string{"main.arch.test"},
		LatestTimestamp: []int64{0},
	}
}
func defRegisArg() *model.ArgRegister {
	return &model.ArgRegister{
		LatestTimestamp: time.Now().Unix(),
		Appid:           "main.arch.test",
		Hostname:        "test1", RPC: "127.0.0.1:8080",
		Region: "shsb", Zone: "sh001",
		Env: "pre", Status: 1,
		Metadata: `{"test":"test","weight":"10"}`,
	}
}

var config = newConfig()

func newConfig() *dc.Config {
	return &dc.Config{HTTPClient: &bm.ClientConfig{Timeout: xtime.Duration(time.Second), Breaker: &breaker.Config{Window: xtime.Duration(time.Second),
		Sleep:   xtime.Duration(time.Millisecond * 100),
		Bucket:  10,
		Ratio:   0.5,
		Request: 100},
		App: &bm.App{Key: "0c4b8fe3ff35a4b6", Secret: "b370880d1aca7d3a289b9b9a7f4d6812"}},
		BM: &dc.HTTPServers{Inner: &bm.ServerConfig{Addr: "127.0.0.1:7171"}},
	}
}

func TestRegister(t *testing.T) {
	Convey("test Register", t, func() {
		svr, _ := New(config)
		i := model.NewInstance(reg)
		svr.Register(context.TODO(), i, reg.LatestTimestamp, false)
		ins, err := svr.Fetch(context.TODO(), fet)
		for _, i := range ins.Instances {
			fmt.Println("ins", i)
		}
		So(err, ShouldBeNil)
		So(len(ins.Instances), ShouldResemble, 1)
		Convey("test metadta", func() {
			for _, i := range ins.Instances {
				So(err, ShouldBeNil)
				So(i.Metadata["weight"], ShouldEqual, "10")
				So(i.Metadata["test"], ShouldEqual, "test")
			}
		})
	})
}
func TestDiscovery(t *testing.T) {
	Convey("test cancel polls", t, func() {
		svr, _ := New(config)
		reg2 := defRegisArg()
		reg2.Hostname = "test2"
		i1 := model.NewInstance(reg)
		i2 := model.NewInstance(reg2)
		svr.Register(context.TODO(), i1, reg.LatestTimestamp, reg.Replication)
		svr.Register(context.TODO(), i2, reg2.LatestTimestamp, reg.Replication)
		ch, new, err := svr.Polls(context.TODO(), pollArg)
		So(err, ShouldBeNil)
		So(new, ShouldBeTrue)
		ins := <-ch
		So(len(ins["main.arch.test"].Instances), ShouldEqual, 2)
		pollArg.LatestTimestamp[0] = ins["main.arch.test"].LatestTimestamp
		time.Sleep(time.Second)
		err = svr.Cancel(context.TODO(), cancel)
		So(err, ShouldBeNil)
		ch, new, err = svr.Polls(context.TODO(), pollArg)
		So(err, ShouldBeNil)
		So(new, ShouldBeTrue)
		ins = <-ch
		So(len(ins["main.arch.test"].Instances), ShouldEqual, 1)
	})
	Convey("test compatible with treeid polls", t, func() {
		svr, _ := New(config)
		reg2 := defRegisArg()
		reg2.Treeid = 1
		i1 := model.NewInstance(reg2)
		svr.Register(ctx, i1, reg2.LatestTimestamp, reg2.Replication)
		ch, new, err := svr.Polls(context.TODO(), pollArg)
		So(err, ShouldBeNil)
		So(new, ShouldBeTrue)
		ins := <-ch
		So(len(ins["main.arch.test"].Instances), ShouldEqual, 1)
		treepoll := newPoll()
		treepoll.Treeid = []int64{1}
		ch, new, err = svr.Polls(context.TODO(), treepoll)
		So(err, ShouldBeNil)
		So(new, ShouldBeTrue)
		ins = <-ch
		So(len(ins["1"].Instances), ShouldEqual, 1)
	})
}
func TestFetchs(t *testing.T) {
	Convey("test fetch multi appid", t, func() {
		svr, _ := New(config)
		reg2 := defRegisArg()
		reg2.Appid = "appid2"
		i1 := model.NewInstance(reg)
		i2 := model.NewInstance(reg2)
		svr.Register(context.TODO(), i1, reg.LatestTimestamp, reg.Replication)
		svr.Register(context.TODO(), i2, reg2.LatestTimestamp, reg.Replication)
		fetchs := newFetchArg()
		fetchs.Appid = append(fetchs.Appid, "appid2")
		is, err := svr.Fetchs(ctx, fetchs)
		So(err, ShouldBeNil)
		So(len(is), ShouldResemble, 2)
	})
}
func TestZones(t *testing.T) {
	Convey("test multi zone discovery", t, func() {
		svr, _ := New(config)
		reg2 := defRegisArg()
		reg2.Zone = "sh002"
		i1 := model.NewInstance(reg)
		i2 := model.NewInstance(reg2)
		svr.Register(context.TODO(), i1, reg.LatestTimestamp, reg.Replication)
		svr.Register(context.TODO(), i2, reg2.LatestTimestamp, reg.Replication)
		ch, new, err := svr.Polls(context.TODO(), pollArg)
		So(err, ShouldBeNil)
		So(new, ShouldBeTrue)
		ins := <-ch
		So(len(ins["main.arch.test"].ZoneInstances), ShouldEqual, 2)
		pollArg.Zone = "sh002"
		ch, new, err = svr.Polls(context.TODO(), pollArg)
		So(err, ShouldBeNil)
		So(new, ShouldBeTrue)
		ins = <-ch
		So(len(ins["main.arch.test"].ZoneInstances), ShouldEqual, 1)
		Convey("test zone update", func() {
			pollArg.LatestTimestamp = []int64{ins["main.arch.test"].LatestTimestamp}
			pollArg.Zone = ""
			reg3 := defRegisArg()
			reg3.Zone = "sh002"
			reg3.Hostname = "test03"
			i3 := model.NewInstance(reg3)
			svr.Register(context.TODO(), i3, reg3.LatestTimestamp, reg3.Replication)
			ch, new, err = svr.Polls(context.TODO(), pollArg)
			So(err, ShouldBeNil)
			ins = <-ch
			So(len(ins["main.arch.test"].ZoneInstances), ShouldResemble, 2)
			So(len(ins["main.arch.test"].ZoneInstances["sh002"]), ShouldResemble, 2)
			So(len(ins["main.arch.test"].ZoneInstances["sh001"]), ShouldResemble, 1)
			pollArg.LatestTimestamp = []int64{ins["main.arch.test"].LatestTimestamp}
			_, _, err = svr.Polls(context.TODO(), pollArg)
			So(err, ShouldResemble, ecode.NotModified)
		})
	})
}
func TestRenew(t *testing.T) {
	Convey("test Renew", t, func() {
		svr, _ := New(config)
		i := model.NewInstance(reg)
		svr.Register(context.TODO(), i, reg.LatestTimestamp, reg.Replication)
		_, err := svr.Renew(context.TODO(), rew)
		So(err, ShouldBeNil)
	})
}

func TestCancel(t *testing.T) {
	Convey("test cancel", t, func() {
		svr, _ := New(config)
		i := model.NewInstance(reg)
		svr.Register(context.TODO(), i, reg.LatestTimestamp, reg.Replication)
		err := svr.Cancel(context.TODO(), cancel)
		So(err, ShouldBeNil)
		_, err = svr.Fetch(context.TODO(), fet)
		So(err, ShouldResemble, ecode.NothingFound)
	})
}

func TestFetchAll(t *testing.T) {
	Convey("test fetch all", t, func() {
		svr, _ := New(config)
		i := model.NewInstance(reg)
		svr.Register(context.TODO(), i, reg.LatestTimestamp, reg.Replication)
		fs := svr.FetchAll(context.TODO())
		_, ok := fs[reg.Appid]
		So(ok, ShouldBeTrue)
	})
}

func TestSet(t *testing.T) {
	Convey("test set", t, func() {
		svr, _ := New(config)
		i := model.NewInstance(reg)
		svr.Register(context.TODO(), i, reg.LatestTimestamp, reg.Replication)
		set.Metadata = []string{`{"weight":"1"}`}
		err := svr.Set(context.TODO(), set)
		So(err, ShouldBeNil)
		cm, err := svr.Fetch(context.TODO(), fet)
		So(err, ShouldBeNil)
		So(cm.Instances[0].Metadata["weight"], ShouldResemble, "1")
	})
}
