package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAccountInfo(t *testing.T) {
	convey.Convey("AccountInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ai, err := d.AccountInfo(c, mid)
			ctx.Convey("Then err should be nil.ai should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ai, convey.ShouldNotBeNil)
			})
		})
	})
}
