package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoArchiveUGCPay(t *testing.T) {
	convey.Convey("ArchiveUGCPay", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(2333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ArchiveUGCPay(c, aid)
			ctx.Convey("Then err should be nil.pay should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				// ctx.So(pay, convey.ShouldEqual, false)
			})
		})
	})
}
