package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDoAvBreach(t *testing.T) {
	convey.Convey("DoAvBreach", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(10)
			aid    = int64(100)
			ctype  = int(1)
			reason = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DoAvBreach(c, mid, aid, ctype, reason)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDoUpForbid(t *testing.T) {
	convey.Convey("DoUpForbid", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(10)
			days   = int(100)
			ctype  = int(1)
			reason = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DoUpForbid(c, mid, days, ctype, reason)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDoUpDismiss(t *testing.T) {
	convey.Convey("DoUpDismiss", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(100)
			ctype  = int(1)
			reason = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DoUpDismiss(c, mid, ctype, reason)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDoUpPass(t *testing.T) {
	convey.Convey("DoUpPass", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mids  = []int64{1, 2, 3}
			ctype = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DoUpPass(c, mids, ctype)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
