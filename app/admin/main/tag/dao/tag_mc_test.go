package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyTag(t *testing.T) {
	var (
		tid = int64(0)
	)
	convey.Convey("keyTag", t, func(ctx convey.C) {
		p1 := keyTag(tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyName(t *testing.T) {
	var (
		name = ""
	)
	convey.Convey("keyName", t, func(ctx convey.C) {
		p1 := keyName(name)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelTagCache(t *testing.T) {
	var (
		c     = context.TODO()
		tid   = int64(0)
		tname = ""
	)
	convey.Convey("DelTagCache", t, func(ctx convey.C) {
		err := d.DelTagCache(c, tid, tname)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
