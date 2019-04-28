package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoHot(t *testing.T) {
	convey.Convey("Hot", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(10100652)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rs, err := d.Hot(c, aid)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}
