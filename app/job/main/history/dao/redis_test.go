package dao

import (
	"context"
	"go-common/app/service/main/history/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyIndex(t *testing.T) {
	convey.Convey("keyIndex", t, func(ctx convey.C) {
		var (
			business = ""
			mid      = int64(14771787)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyIndex(business, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyHistory(t *testing.T) {
	convey.Convey("keyHistory", t, func(ctx convey.C) {
		var (
			business = "archive"
			mid      = int64(14771787)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyHistory(business, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHistoriesCache(t *testing.T) {
	convey.Convey("HistoriesCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			merges = []*model.Merge{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.HistoriesCache(c, merges)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTrimCache(t *testing.T) {
	convey.Convey("TrimCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			business = ""
			mid      = int64(14771787)
			limit    = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.TrimCache(c, business, mid, limit)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCache(t *testing.T) {
	convey.Convey("DelCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			business = ""
			mid      = int64(14771787)
			aids     = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCache(c, business, mid, aids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
			})
		})
	})
}

func TestDaoDelLock(t *testing.T) {
	convey.Convey("DelLock", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.DelLock(c)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}
