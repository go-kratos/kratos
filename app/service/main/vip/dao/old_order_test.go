package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddOldPayOrder(t *testing.T) {
	convey.Convey("AddOldPayOrder", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			orderNo = "20188888888"
			r       = &model.VipOldPayOrder{
				OrderNo: orderNo,
				UserIP:  []byte("127.0.0.1"),
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddOldPayOrder(c, r)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
		convCtx.Convey("clean data", func(convCtx convey.C) {
			d.olddb.Exec(c, "delete from vip_pay_order where order_no=? ", orderNo)
		})
	})
}

func TestDaoAddOldRechargeOrder(t *testing.T) {
	convey.Convey("AddOldRechargeOrder", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			orderNo = "201899999999"
			r       = &model.VipOldRechargeOrder{OrderNo: orderNo}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddOldRechargeOrder(c, r)
			convCtx.So(err, convey.ShouldBeNil)
		})
		convCtx.Convey("clean data", func(convCtx convey.C) {
			d.olddb.Exec(c, "delete from vip_recharge_order where order_no=? ", orderNo)
		})
	})
}
