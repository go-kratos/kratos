package mission

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMissionMissions(t *testing.T) {
	convey.Convey("Missions", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mm, err := d.Missions(c)
			ctx.Convey("Then err should be nil.mm should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mm, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMissionMissionOnlineByTid(t *testing.T) {
	convey.Convey("MissionOnlineByTid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int16(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mm, err := d.MissionOnlineByTid(c, tid)
			ctx.Convey("Then err should be nil.mm should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mm, convey.ShouldNotBeNil)
			})
		})
	})
}
