package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTopPhoto(t *testing.T) {
	convey.Convey("TopPhoto", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			space, err := d.TopPhoto(c, mid)
			ctx.Convey("Then err should be nil.space should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(space, convey.ShouldNotBeNil)
			})
		})
	})
}
