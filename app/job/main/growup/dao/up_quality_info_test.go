package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUpQuality(t *testing.T) {
	convey.Convey("GetUpQuality", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			table = "up_quality_info_11"
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			up, last, err := d.GetUpQuality(c, table, id, limit)
			ctx.Convey("Then err should be nil.up,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(up, convey.ShouldNotBeNil)
			})
		})
	})
}
