package resource

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestResourceVipProducts(t *testing.T) {
	convey.Convey("VipProducts", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := VipProducts(c)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}
