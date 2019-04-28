package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddAuthority(t *testing.T) {
	convey.Convey("AddAuthority", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.AddAuthority(c, mids)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRmAuthority(t *testing.T) {
	convey.Convey("RmAuthority", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.RmAuthority(c, mids)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
