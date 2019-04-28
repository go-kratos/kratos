package upcrm

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmGetUpStatLastDate(t *testing.T) {
	convey.Convey("GetUpStatLastDate", t, func(ctx convey.C) {
		var (
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			lastday, err := d.GetUpStatLastDate(date)
			ctx.Convey("Then err should be nil.lastday should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastday, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmQueryYesterday(t *testing.T) {
	convey.Convey("QueryYesterday", t, func(ctx convey.C) {
		var (
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.QueryYesterday(date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmQueryTrend(t *testing.T) {
	convey.Convey("QueryTrend", t, func(ctx convey.C) {
		var (
			statType    = int(0)
			currentDate = time.Now()
			days        = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.QueryTrend(statType, currentDate, days)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmQueryDetail(t *testing.T) {
	convey.Convey("QueryDetail", t, func(ctx convey.C) {
		var (
			startDate = time.Now()
			endDate   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.QueryDetail(startDate, endDate)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
