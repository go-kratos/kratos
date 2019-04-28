package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawInfo(t *testing.T) {
	convey.Convey("RawInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			r, err := d.RawInfo(c, mid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawInfos(t *testing.T) {
	convey.Convey("RawInfos", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{2089809, 1540883324, 1540818280}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawInfos(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
