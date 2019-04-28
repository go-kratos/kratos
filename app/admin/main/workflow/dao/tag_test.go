package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoChallTagsCount(t *testing.T) {
	convey.Convey("ChallTagsCount", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			gids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			counts, err := d.ChallTagsCount(c, gids)
			ctx.Convey("Then err should be nil.counts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(counts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoChallTagsCountV3(t *testing.T) {
	convey.Convey("ChallTagsCountV3", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			gids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			counts, err := d.ChallTagsCountV3(c, gids)
			ctx.Convey("Then err should be nil.counts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(counts, convey.ShouldNotBeNil)
			})
		})
	})
}
