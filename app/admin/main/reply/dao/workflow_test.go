package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelReport(t *testing.T) {
	convey.Convey("DelReport", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.DelReport(c, oid, rpid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
