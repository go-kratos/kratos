package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddMonitor(t *testing.T) {
	convey.Convey("AddMonitor", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(1)
			operator = "aa"
			remark   = "aa"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMonitor(c, mid, operator, remark)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
func TestDaoMonitors(t *testing.T) {
	convey.Convey("Monitors", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
			pn  = int(0)
			ps  = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mns, total, err := d.Monitors(c, mid, false, pn, ps)
			ctx.Convey("Then err should be nil.mns,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(mns, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelMonitor(t *testing.T) {
	convey.Convey("DelMonitor", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(1)
			operator = "aa"
			remark   = "aa"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelMonitor(c, mid, operator, remark)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
