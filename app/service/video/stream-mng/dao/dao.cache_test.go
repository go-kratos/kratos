package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaostreamFullInfo(t *testing.T) {
	convey.Convey("streamFullInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(11891462)
			sname = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.streamFullInfo(c, id, sname)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaostreamRIDByName(t *testing.T) {
	convey.Convey("streamRIDByName", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.streamRIDByName(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaomultiStreamInfo(t *testing.T) {
	convey.Convey("multiStreamInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{11891462}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.multiStreamInfo(c, keys)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
