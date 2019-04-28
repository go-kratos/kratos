package dao

import (
	"context"
	"go-common/app/admin/main/vip/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoOrderCount(t *testing.T) {
	convey.Convey("OrderCount", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgPayOrder{Mid: 1, OrderNo: "1"}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.OrderCount(c, arg)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoOrderList(t *testing.T) {
	convey.Convey("OrderList", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgPayOrder{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.OrderList(c, arg)
			convCtx.Convey("Then err should be nil.res should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSelOrder(t *testing.T) {
	convey.Convey("SelOrder", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			orderNo = "2016072617212166230921"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.SelOrder(c, orderNo)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddPayOrderLog(t *testing.T) {
	convey.Convey("AddPayOrderLog", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.PayOrderLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddPayOrderLog(c, arg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
