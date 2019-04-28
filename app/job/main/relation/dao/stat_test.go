package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelStatCache(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("DelStatCache", t, func(ctx convey.C) {
		err := d.DelStatCache(mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFollowerAchieve(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(0)
		follower = int64(0)
	)
	convey.Convey("FollowerAchieve", t, func(ctx convey.C) {
		d.FollowerAchieve(c, mid, follower)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestDaoensureAllFollowerAchieve(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(0)
		follower = int64(0)
	)
	convey.Convey("ensureAllFollowerAchieve", t, func(ctx convey.C) {
		d.ensureAllFollowerAchieve(c, mid, follower)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestDaoSendMsg(t *testing.T) {
	var (
		c       = context.Background()
		mid     = int64(0)
		title   = ""
		context = ""
	)
	convey.Convey("SendMsg", t, func(ctx convey.C) {
		err := d.SendMsg(c, mid, title, context)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
