package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetLotteryRIDs(t *testing.T) {
	convey.Convey("GetLotteryRIDs", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			start  = int64(1547436245)
			end    = int64(1547436245)
			offset = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			info, err := d.GetLotteryRIDs(c, start, end, offset)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}
