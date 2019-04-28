package dao

import (
	"context"
	"testing"

	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPayRefund(t *testing.T) {
	convey.Convey("PayRefund", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataJSON = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PayRefund(c, dataJSON)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, ecode.Int(8004020001))
			})
		})
	})
}

func TestDaoPayCancel(t *testing.T) {
	convey.Convey("PayCancel", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataJSON = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PayCancel(c, dataJSON)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, ecode.Int(8004020001))
			})
		})
	})
}

func TestDaoPayQuery(t *testing.T) {
	convey.Convey("PayQuery", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataJSON = `{"customerId":"10017","orderIds":"77546846181122123422","sign":"860d970710eac87650f221d0e0db6940","signType":"MD5","timestamp":"1542882768000","traceId":"1542882768117855000","version":"1.0"}`
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.PayQuery(c, dataJSON)
			ctx.Convey("Then err should be nil.orders should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
