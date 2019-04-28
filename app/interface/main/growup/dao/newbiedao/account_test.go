package newbiedao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewbiedaoGetInfo(t *testing.T) {
	convey.Convey("GetInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515398)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetInfo(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestNewbiedaoGetInfos(t *testing.T) {
	convey.Convey("GetInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{27515398}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetInfos(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
