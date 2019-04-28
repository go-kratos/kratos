package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoequipKey(t *testing.T) {
	convey.Convey("equipKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := equipKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheEquip(t *testing.T) {
	convey.Convey("DelCacheEquip", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheEquip(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
