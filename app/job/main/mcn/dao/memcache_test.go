package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaomcnSignKey(t *testing.T) {
	convey.Convey("mcnSignKey", t, func(ctx convey.C) {
		var (
			mcnMid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mcnSignKey(mcnMid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelMcnSignCache(t *testing.T) {
	convey.Convey("DelMcnSignCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mcnMid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelMcnSignCache(c, mcnMid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
