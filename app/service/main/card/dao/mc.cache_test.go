package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/card/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoequipKey(t *testing.T) {
	convey.Convey("equipKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := equipKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheEquips(t *testing.T) {
	convey.Convey("CacheEquips", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2, 3}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheEquips(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCacheEquip(t *testing.T) {
	convey.Convey("CacheEquip", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.CacheEquip(c, mid)
			ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheEquips(t *testing.T) {
	convey.Convey("AddCacheEquips", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values map[int64]*model.UserEquip
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheEquips(c, values)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheEquip(t *testing.T) {
	convey.Convey("AddCacheEquip", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			v   = &model.UserEquip{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheEquip(c, mid, v)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheEquips(t *testing.T) {
	convey.Convey("DelCacheEquips", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheEquips(c, ids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheEquip(t *testing.T) {
	convey.Convey("DelCacheEquip", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheEquip(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
