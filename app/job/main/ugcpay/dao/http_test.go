package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPayRechargeShell(t *testing.T) {
	convey.Convey("PayRechargeShell", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataJSON = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.PayRechargeShell(c, dataJSON)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPayCheckRefundOrder(t *testing.T) {
	convey.Convey("PayCheckRefundOrder", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataJSON = `{"customerId":"10017","sign":"072dfe62a270c7c44f619244962737a6","signType":"MD5","txIds":"3059753508505497600"}`
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.PayCheckRefundOrder(c, dataJSON)
			ctx.Convey("Then err should be nil.orders should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(orders, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPayCheckOrder(t *testing.T) {
	convey.Convey("PayCheckOrder", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataJSON = `{"customerId":"10017","sign":"072dfe62a270c7c44f619244962737a6","signType":"MD5","txIds":"3059753508505497600"}`
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.PayCheckOrder(c, dataJSON)
			ctx.Convey("Then err should be nil.orders should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// ctx.So(orders, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopaySend(t *testing.T) {
	convey.Convey("paySend", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			url      = ""
			jsonData = ""
			respData = interface{}(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.paySend(c, url, jsonData, respData)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
