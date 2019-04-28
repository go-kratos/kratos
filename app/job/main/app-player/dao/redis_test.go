package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPushList(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("PushList", t, func(ctx convey.C) {
		err := d.PushList(c, nil)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestPopList(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("PopList", t, func(ctx convey.C) {
		_, err := d.PopList(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestPingRedis(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("PingRedis", t, func(ctx convey.C) {
		err := d.PingRedis(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
