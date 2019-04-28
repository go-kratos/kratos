package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeGetArchiveByDate(t *testing.T) {
	convey.Convey("GetArchiveByDate", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			aid   = "av_id"
			table = "av_income"
			date  = "2018-06-24"
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			archives, err := d.GetArchiveByDate(c, aid, table, date, id, limit)
			ctx.Convey("Then err should be nil.archives should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archives, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetBgmIncomeByDate(t *testing.T) {
	convey.Convey("GetBgmIncomeByDate", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = "2018-06-24"
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			archives, err := d.GetBgmIncomeByDate(c, date, id, limit)
			ctx.Convey("Then err should be nil.archives should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archives, convey.ShouldNotBeNil)
			})
		})
	})
}
