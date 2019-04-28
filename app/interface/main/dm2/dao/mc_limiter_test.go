package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyDmDailyLimit(t *testing.T) {
	convey.Convey("keyDmDailyLimit", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyDmDailyLimit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetDmDailyLimitCache(t *testing.T) {
	convey.Convey("SetDmDailyLimitCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			limiter = &model.DailyLimiter{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetDmDailyLimitCache(c, mid, limiter)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetDmDailyLimitCache(t *testing.T) {
	convey.Convey("GetDmDailyLimitCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			limiter, err := testDao.GetDmDailyLimitCache(c, mid)
			ctx.Convey("Then err should be nil.limiter should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(limiter, convey.ShouldNotBeNil)
			})
		})
	})
}
