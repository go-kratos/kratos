package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBangumiPull(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
		ip  = ""
	)
	convey.Convey("BangumiPull", t, func(ctx convey.C) {
		seasonIDS, err := d.BangumiPull(c, mid, ip)
		ctx.Convey("Then err should be nil.seasonIDS should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(seasonIDS, convey.ShouldBeNil)
		})
	})
}

func TestDaoBangumiSeasons(t *testing.T) {
	var (
		c         = context.Background()
		seasonIDs = []int64{5735, 5714, 5702, 5725}
		ip        = ""
	)
	convey.Convey("BangumiSeasons", t, func(ctx convey.C) {
		psm, err := d.BangumiSeasons(c, seasonIDs, ip)
		ctx.Convey("Then err should be nil.psm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(psm, convey.ShouldNotBeNil)
		})
	})
}
