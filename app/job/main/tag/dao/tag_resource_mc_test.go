package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyResTag(t *testing.T) {
	var (
		oid = int64(10002845)
		typ = int32(4)
	)
	convey.Convey("keyResTag", t, func(ctx convey.C) {
		p1 := keyResTag(oid, typ)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelTagResourceCache(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(10002845)
		tp  = int32(4)
	)
	convey.Convey("DelTagResourceCache", t, func(ctx convey.C) {
		err := d.DelTagResourceCache(c, oid, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
