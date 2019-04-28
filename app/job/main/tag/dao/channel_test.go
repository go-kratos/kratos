package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoChannelMap(t *testing.T) {
	convey.Convey("ChannelMap", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tidMap, err := d.ChannelMap(c)
			ctx.Convey("Then err should be nil.tidMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tidMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoChannelRules(t *testing.T) {
	convey.Convey("ChannelRules", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			lastID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rules, err := d.ChannelRules(c, lastID)
			ctx.Convey("Then err should be nil.rules should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rules, convey.ShouldNotBeNil)
			})
		})
	})
}
