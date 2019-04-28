package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyRecent(t *testing.T) {
	convey.Convey("keyRecent", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyRecent(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRecentDM(t *testing.T) {
	convey.Convey("RecentDM", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			start = int64(0)
			end   = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.RecentDM(c, mid, start, end)
		})
	})
}

func TestDaoTrimUpRecent(t *testing.T) {
	convey.Convey("TrimUpRecent", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			count = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.TrimUpRecent(c, mid, count)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
