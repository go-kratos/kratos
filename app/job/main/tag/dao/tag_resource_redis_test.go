package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotagResKey(t *testing.T) {
	var (
		tid = int64(2)
		tp  = int32(0)
	)
	convey.Convey("tagResKey", t, func(ctx convey.C) {
		p1 := tagResKey(tid, tp)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddTagResCache(t *testing.T) {
	var (
		c  = context.TODO()
		rt = &model.ResTag{
			Oid:  1,
			Tids: []int64{1, 2, 3},
		}
	)
	convey.Convey("AddTagResCache", t, func(ctx convey.C) {
		err := d.AddTagResCache(c, rt)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRemoveTagResCache(t *testing.T) {
	var (
		c  = context.TODO()
		rt = &model.ResTag{
			Oid:  1,
			Tids: []int64{1, 2, 3},
		}
	)
	convey.Convey("RemoveTagResCache", t, func(ctx convey.C) {
		err := d.RemoveTagResCache(c, rt)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoresOidKey(t *testing.T) {
	var (
		oid = int64(10002845)
		typ = int32(4)
	)
	convey.Convey("resOidKey", t, func(ctx convey.C) {
		p1 := resOidKey(oid, typ)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelResOidCache(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(10002845)
		tp  = int32(4)
	)
	convey.Convey("DelResOidCache", t, func(ctx convey.C) {
		err := d.DelResOidCache(c, oid, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
