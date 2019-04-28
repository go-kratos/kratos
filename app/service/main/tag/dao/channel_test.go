package dao

import (
	"context"
	"go-common/app/service/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoChannels(t *testing.T) {
	convey.Convey("Channels", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgChannels{
				LastID: 0,
				Size:   50,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Channels(c, arg)
		})
	})
}

func TestDaoChannelRule(t *testing.T) {
	convey.Convey("ChannelRule", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgChannelRule{
				LastID: 0,
				Size:   50,
				State:  0,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ChannelRule(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoChannelCategories(t *testing.T) {
	convey.Convey("ChannelCategories", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgChannelCategory{
				LastID: 0,
				Size:   50,
				State:  0,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ChannelCategories(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoChannelGroup(t *testing.T) {
	convey.Convey("ChannelGroup", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ChannelGroup(c, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
