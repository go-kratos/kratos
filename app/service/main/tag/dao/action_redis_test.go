package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyAction(t *testing.T) {
	convey.Convey("keyAction", t, func(ctx convey.C) {
		var (
			mid = int64(35152246)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyAction(mid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoactionField(t *testing.T) {
	convey.Convey("actionField", t, func(ctx convey.C) {
		var (
			oid = int64(28843596)
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := actionField(oid, tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireAction(t *testing.T) {
	convey.Convey("ExpireAction", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireAction(c, mid, typ)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddActionCache(t *testing.T) {
	convey.Convey("AddActionCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(35152246)
			oid    = int64(28843596)
			tid    = int64(1833)
			typ    = int32(3)
			action = int32(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddActionCache(c, mid, oid, tid, typ, action)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddActionsCache(t *testing.T) {
	convey.Convey("AddActionsCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(35152246)
			typ     = int32(3)
			actions = []*model.ResourceAction{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddActionsCache(c, mid, typ, actions)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoActionCache(t *testing.T) {
	convey.Convey("ActionCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			oid = int64(28843596)
			tid = int64(16906)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			action, err := d.ActionCache(c, mid, oid, tid, typ)
			ctx.Convey("Then err should be nil.action should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(action, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoActionsCache(t *testing.T) {
	convey.Convey("ActionsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(35152246)
			oid  = int64(28843596)
			typ  = int32(3)
			tids = []int64{16906, 1833}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActionsCache(c, mid, oid, typ, tids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelActionCache(t *testing.T) {
	convey.Convey("DelActionCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelActionCache(c, mid, typ)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
