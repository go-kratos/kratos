package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohit(t *testing.T) {
	convey.Convey("hit", t, func(ctx convey.C) {
		var (
			mid = int64(46333)
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			shard := hit(mid)
			ctx.Convey("Then shard should not be mid % 100.", func(ctx convey.C) {
				ctx.So(shard, convey.ShouldEqual, 33)
			})
		})
	})
}

func TestDaoFigureInfo(t *testing.T) {
	convey.Convey("FigureInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(20606508)
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			res, err := d.FigureInfo(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRanks(t *testing.T) {
	convey.Convey("Ranks", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			ranks, err := d.Ranks(c)
			ctx.Convey("Then err should be nil.ranks should have length 100.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ranks, convey.ShouldHaveLength, 100)
			})
		})
	})
}
