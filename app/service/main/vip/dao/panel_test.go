package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoVipPayOrderSuccs(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("VipPayOrderSuccs", t, func(ctx convey.C) {
		mpo, err := d.VipPayOrderSuccs(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("mpo should not be nil", func(ctx convey.C) {
			ctx.So(mpo, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVipPriceConfigs(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("VipPriceConfigs", t, func(ctx convey.C) {
		vpcs, err := d.VipPriceConfigs(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("vpcs should not be nil", func(ctx convey.C) {
			ctx.So(vpcs, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVipPriceDiscountConfigs(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("VipPriceDiscountConfigs", t, func(ctx convey.C) {
		mvp, err := d.VipPriceDiscountConfigs(c)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("mvp should not be nil", func(ctx convey.C) {
			ctx.So(mvp, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVipPriceDiscountByProductID(t *testing.T) {
	var c = context.Background()
	convey.Convey("TestDaoVipPriceDiscountByProductID", t, func(ctx convey.C) {
		mvp, err := d.VipPriceDiscountByProductID(c, "tv.danmaku.bilibilihd.big12month")
		fmt.Println("mvp:", mvp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("mvp should not be nil", func(ctx convey.C) {
			ctx.So(mvp, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVipPriceByProductID(t *testing.T) {
	var c = context.Background()
	convey.Convey("TestDaoVipPriceByProductID", t, func(ctx convey.C) {
		mvp, err := d.VipPriceByProductID(c, "tv.danmaku.bilibilihd.big12month")
		fmt.Println("mvp:", mvp)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("mvp should not be nil", func(ctx convey.C) {
			ctx.So(mvp, convey.ShouldNotBeNil)
		})
	})
}
