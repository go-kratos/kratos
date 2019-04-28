package notice

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRoomNoticeMonthConsume(t *testing.T) {
	convey.Convey("MonthConsume", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(10000)
			targetID = int64(1008)
			date     = "20190101"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MonthConsume(c, id, targetID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
