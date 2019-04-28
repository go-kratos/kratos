package upcrm

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmQueryUpRank(t *testing.T) {
	convey.Convey("QueryUpRank", t, func(ctx convey.C) {
		var (
			rankType = int(0)
			date     = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryUpRank(rankType, date)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmQueryUpRankAll(t *testing.T) {
	convey.Convey("QueryUpRankAll", t, func(ctx convey.C) {
		var (
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryUpRankAll(date)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetUpRankLatestDate(t *testing.T) {
	convey.Convey("GetUpRankLatestDate", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			date, err := d.GetUpRankLatestDate()
			ctx.Convey("Then err should be nil.date should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(date, convey.ShouldNotBeNil)
			})
		})
	})
}
