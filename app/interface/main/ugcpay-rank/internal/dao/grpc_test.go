package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoElecMonthRank(t *testing.T) {
	convey.Convey("ElecMonthRank", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			avID     = int64(0)
			upMID    = int64(0)
			rankSize = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			resp, err := d.RankElecMonth(c, avID, upMID, rankSize)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(resp, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoElecMonthRankUP(t *testing.T) {
	convey.Convey("ElecMonthRankUP", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			upMID    = int64(0)
			rankSize = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.RankElecMonthUP(c, upMID, rankSize)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoElecTotalRankAV(t *testing.T) {
	convey.Convey("ElecTotalRankAV", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			upMID    = int64(0)
			avID     = int64(0)
			rankSize = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rank, err := d.RankElecAllAV(c, upMID, avID, rankSize)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}
