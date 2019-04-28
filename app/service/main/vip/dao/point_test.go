package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAllPointExchangePrice(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("AllPointExchangePrice", t, func(ctx convey.C) {
		pe, err := d.AllPointExchangePrice(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("pe should not be nil", func(ctx convey.C) {
			ctx.So(pe, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoPointExchangePrice(t *testing.T) {
	var (
		c     = context.TODO()
		month = int16(0)
	)
	convey.Convey("PointExchangePrice", t, func(ctx convey.C) {
		_, err := d.PointExchangePrice(c, month)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
