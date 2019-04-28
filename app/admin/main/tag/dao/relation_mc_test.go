package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyResource(t *testing.T) {
	var (
		oid = int64(0)
		tid = int64(0)
		typ = int32(0)
	)
	convey.Convey("keyResource", t, func(ctx convey.C) {
		p1 := keyResource(oid, tid, typ)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelResMemCache(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(0)
		tid = int64(0)
		tp  = int32(0)
	)
	convey.Convey("DelResMemCache", t, func(ctx convey.C) {
		err := d.DelResMemCache(c, oid, tid, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
