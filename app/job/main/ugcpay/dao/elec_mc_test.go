package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddCacheOrderID(t *testing.T) {
	convey.Convey("AddCacheOrderID", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "233"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.AddCacheOrderID(c, orderID)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheOrderID(t *testing.T) {
	convey.Convey("DelCacheOrderID", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "233"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheOrderID(c, orderID)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
