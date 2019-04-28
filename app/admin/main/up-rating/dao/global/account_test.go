package global

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGlobalMID(t *testing.T) {
	convey.Convey("MID", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			nickname = "hashbaz"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mid, err := MID(c, nickname)
			ctx.Convey("Then err should be nil.mid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGlobalNames(t *testing.T) {
	convey.Convey("Names", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{40052}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := Names(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
