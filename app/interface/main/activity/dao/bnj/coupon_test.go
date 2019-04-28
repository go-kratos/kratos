package bnj

import (
	"context"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestBnjGrantCoupon(t *testing.T) {
	convey.Convey("GrantCoupon", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(2080809)
			couponID = "3d005e8ba01c5cb0"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.grantCouponURL).Reply(200).JSON(`{"code":0}`)
			err := d.GrantCoupon(c, mid, couponID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
