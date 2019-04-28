package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPushMsg(t *testing.T) {
	var (
		c           = context.Background()
		op          = int32(0)
		server      = ""
		msg         = ""
		keys        = []string{"key"}
		contentType = int32(0)
	)
	convey.Convey("PushMsg", t, func(ctx convey.C) {
		err := d.PushMsg(c, op, server, msg, keys, contentType)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBroadcastRoomMsg(t *testing.T) {
	var (
		c           = context.Background()
		op          = int32(0)
		room        = ""
		msg         = ""
		contentType = int32(0)
	)
	convey.Convey("BroadcastRoomMsg", t, func(ctx convey.C) {
		err := d.BroadcastRoomMsg(c, op, room, msg, contentType)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBroadcastMsg(t *testing.T) {
	var (
		c           = context.Background()
		op          = int32(0)
		speed       = int32(0)
		msg         = ""
		platform    = ""
		contentType = int32(0)
	)
	convey.Convey("BroadcastMsg", t, func(ctx convey.C) {
		err := d.BroadcastMsg(c, op, speed, msg, platform, contentType)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
