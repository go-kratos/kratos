package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRegions(t *testing.T) {
	convey.Convey("Regions", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Regions(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				println(len(res))
			})
		})
	})
}

func TestDaoFindLastMtime(t *testing.T) {
	convey.Convey("Regions", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Regions(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				println(len(res))
			})
		})
	})
}
