package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosubKey(t *testing.T) {
	convey.Convey("subKey", t, func(ctx convey.C) {
		var (
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.subKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireSubCache(t *testing.T) {
	convey.Convey("ExpireSubCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireSubCache(c, mid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddSubMapCache(t *testing.T) {
	convey.Convey("AddSubMapCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			subs = map[int64]*model.Sub{
				1833: {
					Mid: 35152246,
					Tid: 1833,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddSubMapCache(c, mid, subs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddSubListCache(t *testing.T) {
	convey.Convey("AddSubListCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			subs = []*model.Sub{
				{
					Mid: 35152246,
					Tid: 1833,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddSubListCache(c, mid, subs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelSubTidCache(t *testing.T) {
	convey.Convey("DelSubTidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelSubTidCache(c, mid, tid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelSubTidsCache(t *testing.T) {
	convey.Convey("DelSubTidsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelSubTidsCache(c, mid, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelSubCache(t *testing.T) {
	convey.Convey("DelSubCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelSubCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoIsSubCache(t *testing.T) {
	convey.Convey("IsSubCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.IsSubCache(c, mid, tid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsSubsCache(t *testing.T) {
	convey.Convey("IsSubsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.IsSubsCache(c, mid, tids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubTidsCache(t *testing.T) {
	convey.Convey("SubTidsCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SubTidsCache(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubCache(t *testing.T) {
	convey.Convey("SubCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			subs, rem, err := d.SubCache(c, mid)
			ctx.Convey("Then err should be nil.subs,rem should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rem, convey.ShouldNotBeNil)
				ctx.So(subs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubListCache(t *testing.T) {
	convey.Convey("SubListCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			subs, err := d.SubListCache(c, mid)
			ctx.Convey("Then err should be nil.subs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(subs, convey.ShouldNotBeNil)
			})
		})
	})
}
