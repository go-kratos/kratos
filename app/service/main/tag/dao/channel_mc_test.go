package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyChannelGroup(t *testing.T) {
	var (
		tid = int64(1)
	)
	convey.Convey("keyChannelGroup", t, func(ctx convey.C) {
		p1 := keyChannelGroup(tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddChannelGroupCache(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(1)
		cgs = []*model.ChannelGroup{}
	)
	convey.Convey("AddChannelGroupCache", t, func(ctx convey.C) {
		err := d.AddChannelGroupCache(c, tid, cgs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoChannelGroupCache(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(1)
	)
	convey.Convey("ChannelGroupCache", t, func(ctx convey.C) {
		res, err := d.ChannelGroupCache(c, tid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}
