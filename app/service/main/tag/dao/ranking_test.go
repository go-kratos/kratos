package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRankHots(t *testing.T) {
	convey.Convey("RankHots", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.RankHots(c)
		})
	})
}

func TestDaoBangumis(t *testing.T) {
	convey.Convey("Bangumis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Bangumis(c)
		})
	})
}

func TestDaoRegions(t *testing.T) {
	convey.Convey("Regions", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(9)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ps, err := d.Regions(c, rid)
			ctx.Convey("Then err should be nil.ps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ps, convey.ShouldNotBeNil)
			})
		})
	})
}
