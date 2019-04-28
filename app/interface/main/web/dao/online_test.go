package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoOnlineCount(t *testing.T) {
	convey.Convey("OnlineCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.OnlineCount(c)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLiveOnlineCount(t *testing.T) {
	convey.Convey("LiveOnlineCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.LiveOnlineCount(c)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOnlineList(t *testing.T) {
	convey.Convey("OnlineList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			num = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.OnlineList(c, num)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}
