package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoChannelCategories(t *testing.T) {
	var (
		c        = context.TODO()
		lastID   = int64(0)
		pageSize = int32(0)
		state    = int32(0)
	)
	convey.Convey("ChannelCategories", t, func(ctx convey.C) {
		res, err := d.ChannelCategories(c, lastID, pageSize, state)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoChannels(t *testing.T) {
	var (
		c        = context.TODO()
		lastID   = int64(0)
		pageSize = int32(0)
	)
	convey.Convey("Channels", t, func(ctx convey.C) {
		res, err := d.Channels(c, lastID, pageSize)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoChannelRules(t *testing.T) {
	var (
		c        = context.TODO()
		lastID   = int64(0)
		pageSize = int32(0)
		state    = int32(0)
	)
	convey.Convey("ChannelRules", t, func(ctx convey.C) {
		res, err := d.ChannelRules(c, lastID, pageSize, state)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldHaveLength, 0)
		})
	})
}

func TestDaoChannelGroup(t *testing.T) {
	var (
		c   = context.TODO()
		tid = int64(2)
	)
	convey.Convey("ChannelGroup", t, func(ctx convey.C) {
		res, err := d.ChannelGroup(c, tid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
