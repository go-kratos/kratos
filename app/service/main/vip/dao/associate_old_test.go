package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoActivityOrder(t *testing.T) {
	convey.Convey("ActivityOrder", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			orderNo = "test_activity_no"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.ActivityOrder(c, orderNo)
			convCtx.So(err, convey.ShouldBeNil)
		})
		convCtx.Convey("clean data", func(convCtx convey.C) {
			d.olddb.Exec(c, "delete from vip_order_activity_record where order_no=? ", orderNo)
		})
	})
}

func TestDaoUpdateActivityState(t *testing.T) {
	convey.Convey("UpdateActivityState", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			state   = int8(0)
			orderNO = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			aff, err := d.UpdateActivityState(c, state, orderNO)
			convCtx.Convey("Then err should be nil.aff should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(aff, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountProductBuy(t *testing.T) {
	convey.Convey("CountProductBuy", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(0)
			months    = int32(0)
			panelType = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CountProductBuy(c, mid, months, panelType)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
