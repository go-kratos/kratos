package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAreaList(t *testing.T) {
	convey.Convey("AreaList", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything looks good.", func(ctx convey.C) {
			list, err := d.AreaList(c)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func Test_AreaLastTime(t *testing.T) {
	convey.Convey("AreaLastTime", t, func(ctx convey.C) {
		res, err := d.AreaLastTime(context.Background())
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(res, convey.ShouldBeGreaterThan, 0)
	})
}
