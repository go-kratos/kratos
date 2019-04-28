package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaooidKey(t *testing.T) {
	convey.Convey("oidKey", t, func(ctx convey.C) {
		var (
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := oidKey(oid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireOidCache(t *testing.T) {
	convey.Convey("ExpireOidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireOidCache(c, oid, typ)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOidCache(t *testing.T) {
	convey.Convey("OidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.OidCache(c, oid, typ)
		})
	})
}

func TestDaoAddOidCache(t *testing.T) {
	convey.Convey("AddOidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
			r   = &model.Resource{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddOidCache(c, oid, typ, r)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddOidsCache(t *testing.T) {
	convey.Convey("AddOidsCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
			rs  = []*model.Resource{
				{
					Oid:  28843596,
					Type: 3,
					Tid:  1833,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddOidsCache(c, oid, typ, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddOidMapCache(t *testing.T) {
	convey.Convey("AddOidMapCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(28843596)
			typ   = int32(3)
			rsMap = map[int64]*model.Resource{
				28843596: {
					Oid:  28843596,
					Type: 3,
					Tid:  1833,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddOidMapCache(c, oid, typ, rsMap)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoZremOidCache(t *testing.T) {
	convey.Convey("ZremOidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			tid = int64(1833)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ZremOidCache(c, oid, tid, typ)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelOidCache(t *testing.T) {
	convey.Convey("DelOidCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelOidCache(c, oid, typ)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
