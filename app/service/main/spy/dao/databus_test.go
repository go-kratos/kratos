package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPubScoreChange(t *testing.T) {
	convey.Convey("PubScoreChange", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
			msg = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PubScoreChange(c, mid, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
