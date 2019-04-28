package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUnicomGiftState(t *testing.T) {
	convey.Convey("UnicomGiftState", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515247)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			state, err := d.UnicomGiftState(c, mid)
			ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}
