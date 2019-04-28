package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPriceConfigsByStatus(t *testing.T) {
	convey.Convey("PriceConfigsByStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			status = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, total, err := d.PriceConfigsByStatus(c, status)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSaledPriceConfigsByStatus(t *testing.T) {
	convey.Convey("SaledPriceConfigsByStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			status = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, total, err := d.SaledPriceConfigsByStatus(c, status)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
