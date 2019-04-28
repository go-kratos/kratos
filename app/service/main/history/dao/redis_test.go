package dao

import (
	"context"
	"testing"

	pb "go-common/app/service/main/history/api/grpc"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoListCache(t *testing.T) {
	var (
		c        = context.Background()
		business = "pgc"
		mid      = int64(1)
		start    = int64(0)
		his      = []*pb.AddHistoryReq{{
			Mid:      1,
			Business: "pgc",
			Kid:      1,
			Aid:      2,
			Sid:      3,
		},
		}
		h = &pb.AddHistoryReq{
			Mid:      2,
			Business: "pgc",
			Kid:      1,
			Aid:      2,
			Sid:      3,
		}
	)
	convey.Convey("add his", t, func() {
		convey.So(d.AddHistoriesCache(c, his), convey.ShouldBeNil)
		convey.So(d.AddHistoryCache(c, h), convey.ShouldBeNil)
		convey.Convey("ListCacheByTime", func(ctx convey.C) {
			aids, err := d.ListCacheByTime(c, business, mid, start)
			ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(aids, convey.ShouldNotBeEmpty)
			})
		})
		convey.Convey("ListsCacheByTime", func(ctx convey.C) {
			var (
				c          = context.Background()
				businesses = []string{"pgc"}
				viewAt     = int64(100)
				ps         = int64(1)
			)
			res, err := d.ListsCacheByTime(c, businesses, mid, viewAt, ps)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})

		convey.Convey("HistoriesCache", func(ctx convey.C) {
			var hs = map[string][]int64{"pgc": {1}}
			res, err := d.HistoriesCache(c, 2, hs)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})

		convey.Convey("ClearHistoryCache", func(ctx convey.C) {
			err := d.ClearHistoryCache(c, mid, []string{business})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		convey.Convey("DelHistoryCache", func(ctx convey.C) {
			ctx.So(d.DelHistoryCache(c, &pb.DelHistoriesReq{
				Mid: 1, Records: []*pb.DelHistoriesReq_Record{{ID: 1, Business: "pgc"}},
			}), convey.ShouldBeNil)
		})
		convey.Convey("TrimCache", func(ctx convey.C) {
			err := d.TrimCache(c, business, mid, 10)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})

	})
}

func TestDaoDelCache(t *testing.T) {
	var (
		c        = context.Background()
		business = "pgc"
		mid      = int64(1)
		aids     = []int64{1}
	)
	convey.Convey("DelCache", t, func(ctx convey.C) {
		err := d.DelCache(c, business, mid, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetUserHideCache(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(1)
		value = int64(1)
	)
	convey.Convey("SetUserHideCache", t, func(ctx convey.C) {
		err := d.SetUserHideCache(c, mid, value)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUserHideCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("UserHideCache", t, func(ctx convey.C) {
		value, err := d.UserHideCache(c, mid)
		ctx.Convey("Then err should be nil.value should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(value, convey.ShouldNotBeNil)
			// ctx.So(value, convey.ShouldEqual, 1)
		})
	})
}
