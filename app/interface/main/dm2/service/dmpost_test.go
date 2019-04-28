package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenDMID(t *testing.T) {
	dm := &model.DM{
		Progress: 0,
		Pool:     0,
		State:    0,
		Mid:      1001,
		Content: &model.Content{
			Msg:      "dm msg test",
			Mode:     1,
			FontSize: 25,
			Color:    165000,
			IP:       100097834,
		},
	}
	time.Sleep(time.Second * 1)
	Convey("err should be nil", t, func() {
		err := svr.genDMID(context.Background(), dm)
		t.Logf("dmid:%v", dm.ID)
		So(err, ShouldBeNil)
	})
}

func TestPost(t *testing.T) {
	dm := &model.DM{
		ID:       1234,
		Type:     1,
		Oid:      1221,
		Progress: 0,
		Pool:     0,
		State:    0,
		Mid:      1001,
		Content: &model.Content{
			ID:       1234,
			Msg:      "dm msg test",
			Mode:     1,
			FontSize: 25,
			Color:    165000,
			IP:       100097834,
		},
	}
	time.Sleep(time.Second * 1)
	Convey("check dm post", t, func() {
		err := svr.Post(context.TODO(), dm, 1234, 456)
		So(err, ShouldBeNil)
	})
}

func TestCheckOversea(t *testing.T) {
	time.Sleep(time.Second * 1)
	Convey("check oversea user", t, func() {
		err := svr.checkOverseasUser(context.TODO())
		So(err, ShouldEqual, ecode.ServiceUpdate)
	})
}

func TestCheckMonitor(t *testing.T) {
	dm := &model.DM{
		ID:       1234,
		Oid:      1221,
		Progress: 0,
		Pool:     0,
		State:    0,
		Mid:      1001,
		Content: &model.Content{
			ID:       1234,
			Msg:      "dm msg test",
			Mode:     1,
			FontSize: 25,
			Color:    165000,
			IP:       100097834,
		},
	}
	sub := &model.Subject{
		Type:  1,
		Oid:   1221,
		State: 0,
		Attr:  16,
	}
	Convey("check monitor type before", t, func() {
		err := svr.checkMonitor(context.TODO(), sub, dm)
		So(err, ShouldBeNil)
		So(dm.State, ShouldEqual, model.StateMonitorBefore)
	})

	dm.State = 0
	sub.Attr = 32
	Convey("check monitor type after", t, func() {
		err := svr.checkMonitor(context.TODO(), sub, dm)
		So(err, ShouldBeNil)
		So(dm.State, ShouldEqual, model.StateMonitorAfter)
	})

}

func TestCheckUpFilter(t *testing.T) {
	dm := &model.DM{
		ID:       1234,
		Oid:      1,
		Progress: 0,
		Pool:     0,
		State:    0,
		Mid:      1001,
		Content: &model.Content{
			ID:       1234,
			Msg:      `[0,0,"1-1",4.5,"adsfasdf",0,0,0,0,500,0,true,"黑体",1]`,
			Mode:     7,
			FontSize: 25,
			Color:    165000,
			IP:       100097834,
		},
	}
	Convey("test up filter", t, func() {
		svr.checkUpFilter(context.TODO(), dm, 1, 2)
	})
}
