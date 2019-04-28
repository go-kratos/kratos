package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeySubSort(t *testing.T) {
	convey.Convey("keySubSort", t, func(ctx convey.C) {
		var (
			mid = int64(35152246)
			typ = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keySubSort(mid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubSortCache(t *testing.T) {
	convey.Convey("SubSortCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			typ = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.SubSortCache(c, mid, typ)
		})
	})
}

func TestDaoAddSubSortCache(t *testing.T) {
	convey.Convey("AddSubSortCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			typ  = int(0)
			tids = []int64{1822, 1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddSubSortCache(c, mid, typ, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelSubSortCache(t *testing.T) {
	convey.Convey("DelSubSortCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			typ = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelSubSortCache(c, mid, typ)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
