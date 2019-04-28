package kfc

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestKfcKfcCoupon(t *testing.T) {
	convey.Convey("KfcCoupon", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(3)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.KfcCoupon(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
