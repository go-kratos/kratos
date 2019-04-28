package pendant

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPendantVipInfo(t *testing.T) {
	convey.Convey("VipInfo", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			idt, err := d.VipInfo(c, mid, ip)
			ctx.Convey("Then err should be nil.idt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(idt, convey.ShouldNotBeNil)
			})
		})
	})
}
