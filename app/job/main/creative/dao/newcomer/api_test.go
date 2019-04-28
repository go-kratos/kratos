package newcomer

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewcomerSendNotify(t *testing.T) {
	convey.Convey("Msg", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mids    = []int64{27515405}
			mc      = "1_17_2"
			title   = "creative"
			context = "sssss"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendNotify(c, mids, mc, title, context)
			ctx.Convey("Then err should be nil.msg should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}
