package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohashRowKey(t *testing.T) {
	convey.Convey("hashRowKey", t, func(ctx convey.C) {
		var (
			tid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hashRowKey(tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWeightLog(t *testing.T) {
	convey.Convey("WeightLog", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			taskid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ls, err := d.WeightLog(c, taskid)
			ctx.Convey("Then err should be nil.ls should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ls, convey.ShouldNotBeNil)
			})
		})
	})
}
