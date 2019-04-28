package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoEquips(t *testing.T) {
	convey.Convey("Equips", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{1, 2, 3, 4}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Equips(c, keys)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoEquip(t *testing.T) {
	convey.Convey("Equip", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Equip(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
