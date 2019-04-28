package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLimitRes(t *testing.T) {
	convey.Convey("LimitRes", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tye = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LimitRes(c, tye)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWhiteUser(t *testing.T) {
	convey.Convey("WhiteUser", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			midm, err := d.WhiteUser(c)
			ctx.Convey("Then err should be nil.midm should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(midm, convey.ShouldNotBeNil)
			})
		})
	})
}
