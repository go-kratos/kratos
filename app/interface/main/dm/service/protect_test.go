package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddProtectApply(t *testing.T) {
	Convey("test add  protect apply", t, func() {
		err := svr.AddProtectApply(context.TODO(), 27515256, 10108163, []int64{719925843, 719925844})
		So(err, ShouldBeNil)
	})
}

func TestUptPaSwitch(t *testing.T) {
	Convey("test upt  pa switch", t, func() {
		err := svr.UptPaSwitch(context.TODO(), 27515256, 1)
		So(err, ShouldBeNil)
		err = svr.UptPaSwitch(context.TODO(), 27515256, 0)
		So(err, ShouldBeNil)
	})
}

func TestUptPaStatus(t *testing.T) {
	Convey("test upt pa status", t, func() {
		err := svr.UptPaStatus(context.TODO(), 27515256, []int64{541, 542}, 1)
		So(err, ShouldBeNil)
		err = svr.UptPaStatus(context.TODO(), 27515256, []int64{541, 542}, -1)
		So(err, ShouldBeNil)
	})
}

func TestProtectApplies(t *testing.T) {
	Convey("test protect applies", t, func() {
		res, err := svr.ProtectApplies(context.TODO(), 27515256, 10097377, 1, "playtime")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestPaVideoLs(t *testing.T) {
	Convey("test pa video ls", t, func() {
		_, err := svr.PaVideoLs(context.TODO(), 27515256)
		So(err, ShouldBeNil)
	})
}
