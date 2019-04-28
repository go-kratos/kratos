package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoFilterAreas(t *testing.T) {
	convey.Convey("FilterAreas", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			source = int64(0)
			area   = "common"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			fs, err := d.FilterAreas(c, source, area)
			ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fs, convey.ShouldNotBeNil)
			})
		})
	})
}
