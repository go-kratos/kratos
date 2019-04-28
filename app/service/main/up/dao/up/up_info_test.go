package up

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpRawUpInfoActive(t *testing.T) {
	convey.Convey("RawUpInfoActive", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upInfoActive, err := d.RawUpInfoActive(c, mid)
			ctx.Convey("Then err should be nil.upInfoActive should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upInfoActive, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpRawUpsInfoActive(t *testing.T) {
	convey.Convey("RawUpsInfoActive", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawUpsInfoActive(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
