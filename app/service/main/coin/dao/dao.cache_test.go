package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUserCoin(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(123)
	)
	convey.Convey("UserCoin", t, func(ctx convey.C) {
		res, err := d.UserCoin(c, id)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoItemCoin(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(123)
		tp = int64(60)
	)
	convey.Convey("ItemCoin", t, func(ctx convey.C) {
		res, err := d.ItemCoin(c, id, tp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
