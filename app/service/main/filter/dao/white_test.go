package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoWhiteAreas(t *testing.T) {
	convey.Convey("WhiteAreas", t, func(ctx convey.C) {
		var (
			c    = context.TODO()
			area = "common"
		)
		ctx.Convey("When everything looks good.", func(ctx convey.C) {
			fs, err := d.WhiteAreas(c, area)
			ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fs, convey.ShouldNotBeNil)
			})
		})
	})
}
