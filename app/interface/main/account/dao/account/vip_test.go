package account

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAccountVipPointBalance(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("VipPointBalance", t, func(ctx convey.C) {
		pointBalance, err := d.VipPointBalance(c, mid, ip)
		ctx.Convey("Then err should be nil.pointBalance should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(pointBalance, convey.ShouldNotBeNil)
		})
	})
}
