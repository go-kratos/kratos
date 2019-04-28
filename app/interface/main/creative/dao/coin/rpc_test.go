package coin

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCoinAddCoin(t *testing.T) {
	convey.Convey("AddCoin", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(27515256)
			aid  = int64(10110817)
			coin = float64(-2.0)
			ip   = "127.0.0.1"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCoin(c, mid, aid, coin, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
