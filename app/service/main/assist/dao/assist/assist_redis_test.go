package assist

import (
	"context"
	"go-common/app/service/main/assist/model/assist"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAssistlogCountKey(t *testing.T) {
	var (
		mid       = int64(0)
		assistMid = int64(0)
	)
	convey.Convey("logCountKey", t, func(ctx convey.C) {
		p1 := d.logCountKey(mid, assistMid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistassTotalKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("assTotalKey", t, func(ctx convey.C) {
		p1 := d.assTotalKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistassSameKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("assSameKey", t, func(ctx convey.C) {
		p1 := d.assSameKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistassUpKey(t *testing.T) {
	var (
		assistMid = int64(0)
	)
	convey.Convey("assUpKey", t, func(ctx convey.C) {
		p1 := d.assUpKey(assistMid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistDailyLogCount(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
		tp        = int64(0)
	)
	convey.Convey("DailyLogCount", t, func(ctx convey.C) {
		count, err := d.DailyLogCount(c, mid, assistMid, tp)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistIncrLogCount(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
		tp        = int64(0)
	)
	convey.Convey("IncrLogCount", t, func(ctx convey.C) {
		err := d.IncrLogCount(c, mid, assistMid, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistTotalAssCnt(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("TotalAssCnt", t, func(ctx convey.C) {
		count, err := d.TotalAssCnt(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistSameAssCnt(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
	)
	convey.Convey("SameAssCnt", t, func(ctx convey.C) {
		count, err := d.SameAssCnt(c, mid, assistMid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistIncrAssCnt(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(0)
		assistMid = int64(0)
	)
	convey.Convey("IncrAssCnt", t, func(ctx convey.C) {
		err := d.IncrAssCnt(c, mid, assistMid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistDelAssUpAllCache(t *testing.T) {
	var (
		c         = context.Background()
		assistMid = int64(0)
	)
	convey.Convey("DelAssUpAllCache", t, func(ctx convey.C) {
		err := d.DelAssUpAllCache(c, assistMid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistAddAssUpAllCache(t *testing.T) {
	var (
		c         = context.Background()
		assistMid = int64(0)
		ups       map[int64]*assist.Up
	)
	convey.Convey("AddAssUpAllCache", t, func(ctx convey.C) {
		err := d.AddAssUpAllCache(c, assistMid, ups)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistAssUpCacheWithScore(t *testing.T) {
	var (
		c         = context.Background()
		assistMid = int64(0)
		start     = int64(0)
		end       = int64(0)
	)
	convey.Convey("AssUpCacheWithScore", t, func(ctx convey.C) {
		_, _, _, err := d.AssUpCacheWithScore(c, assistMid, start, end)
		ctx.Convey("Then err should be nil.mids,ups,total should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
