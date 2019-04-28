package dao

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawStreamFullInfo(t *testing.T) {
	convey.Convey("RawStreamFullInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(11891462)
			sname = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawStreamFullInfo(c, id, sname)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawStreamRIDByName(t *testing.T) {
	convey.Convey("RawStreamRIDByName", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sname = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawStreamRIDByName(c, sname)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawMultiStreamInfo(t *testing.T) {
	convey.Convey("RawMultiStreamInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rids = []int64{11891462}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawMultiStreamInfo(c, rids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaodiffStreamInfo(t *testing.T) {
	convey.Convey("diffStreamInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			info = &model.StreamFullInfo{}
			main = &model.MainStream{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.diffStreamInfo(c, info, main)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
