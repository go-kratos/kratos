package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	day = "20180808"
	mid = int64(91221505)
	app = int64(1)
)

func incrLimitDayCache() (int, error) {
	return d.IncrLimitDayCache(context.Background(), day, app, mid)
}

func TestDaoLimitDayCache(t *testing.T) {
	incrLimitDayCache()
	convey.Convey("LimitDayCache", t, func(ctx convey.C) {
		count, err := d.LimitDayCache(context.Background(), day, app, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldBeGreaterThan, 0)
		})
	})
}

func TestDaoIncrLimitDayCache(t *testing.T) {
	convey.Convey("IncrLimitDayCache", t, func(ctx convey.C) {
		count, err := incrLimitDayCache()
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldBeGreaterThan, 0)
		})
	})
}

func incrLimitNotLiveCache() (int, error) {
	return d.IncrLimitNotLiveCache(context.Background(), day, mid)
}

func TestDaoIncrLimitBizCache(t *testing.T) {
	incrLimitNotLiveCache()
	convey.Convey("IncrLimitBizCache", t, func(ctx convey.C) {
		count, err := d.IncrLimitBizCache(context.Background(), day, app, mid, biz)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldBeGreaterThan, 0)
		})
	})
}

func TestDaoIncrLimitNotLiveCache(t *testing.T) {
	convey.Convey("IncrLimitNotLiveCache", t, func(ctx convey.C) {
		count, err := incrLimitNotLiveCache()
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldBeGreaterThan, 0)
		})
	})
}
