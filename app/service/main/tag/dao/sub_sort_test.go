package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddSubSort(t *testing.T) {
	convey.Convey("AddSubSort", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			typ  = int(3)
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := d.AddSubSort(c, mid, typ, tids)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddSubChannel(t *testing.T) {
	convey.Convey("AddSubChannel", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			tp   = int(3)
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddSubChannel(c, mid, tp, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddSubChannels(t *testing.T) {
	convey.Convey("AddSubChannels", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(35152246)
			tidMap = map[int32][]int64{
				0: {1833, 12096},
				1: {1833, 12096},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddSubChannels(c, mid, tidMap)
			ctx.Convey("TestDaoAddSubChannels Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCustomSubSort(t *testing.T) {
	convey.Convey("CustomSubSort", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			typ = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tids, err := d.CustomSubSort(c, mid, typ)
			ctx.Convey("Then err should be nil.tids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllCustomSubSort(t *testing.T) {
	convey.Convey("AllCustomSubSort", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tids, err := d.AllCustomSubSort(c, mid)
			ctx.Convey("AllCustomSubSort Then err should be nil.tids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tids, convey.ShouldNotBeNil)
			})
		})
	})
}
