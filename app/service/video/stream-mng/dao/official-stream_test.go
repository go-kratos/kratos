package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetOfficialStreamByName(t *testing.T) {
	convey.Convey("GetOfficialStreamByName", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetOfficialStreamByName(c, name)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetOfficialStreamByRoomID(t *testing.T) {
	convey.Convey("GetOfficialStreamByRoomID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetOfficialStreamByRoomID(c, rid)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}
