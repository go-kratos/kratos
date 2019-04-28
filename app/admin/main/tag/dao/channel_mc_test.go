package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyChannelGroup(t *testing.T) {
	var (
		tid = int64(0)
	)
	convey.Convey("keyChannelGroup", t, func(ctx convey.C) {
		p1 := keyChannelGroup(tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelChannelGroupCache(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(0)
	)
	convey.Convey("DelChannelGroupCache", t, func(ctx convey.C) {
		err := d.DelChannelGroupCache(c, tid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
