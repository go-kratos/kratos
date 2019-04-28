package upcrmservice

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmserviceQueryYesterday(t *testing.T) {
	convey.Convey("QueryYesterday", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			item, err := s.QueryYesterday(c, date)
			ctx.Convey("Then err should be nil.item should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(item, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceQueryTrend(t *testing.T) {
	convey.Convey("QueryTrend", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			statType    = int(0)
			currentDate = time.Now()
			days        = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.QueryTrend(c, statType, currentDate, days)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceQueryDetail(t *testing.T) {
	convey.Convey("QueryDetail", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			currentDate = time.Now()
			days        = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.QueryDetail(c, currentDate, days)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
