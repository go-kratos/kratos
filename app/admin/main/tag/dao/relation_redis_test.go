package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotagResKey(t *testing.T) {
	var (
		tid = int64(0)
		typ = int32(0)
	)
	convey.Convey("tagResKey", t, func(ctx convey.C) {
		p1 := tagResKey(tid, typ)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaooidKey(t *testing.T) {
	var (
		oid = int64(0)
		typ = int32(0)
	)
	convey.Convey("oidKey", t, func(ctx convey.C) {
		p1 := oidKey(oid, typ)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelRelationCache(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(0)
		tid = int64(0)
		typ = int32(0)
	)
	convey.Convey("DelRelationCache", t, func(ctx convey.C) {
		err := d.DelRelationCache(c, oid, tid, typ)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
