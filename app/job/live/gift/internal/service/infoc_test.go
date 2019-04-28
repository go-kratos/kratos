package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicebagLogInfoc(t *testing.T) {
	convey.Convey("bagLogInfoc", t, func(ctx convey.C) {
		var (
			uid      = int64(0)
			bagID    = int64(0)
			giftID   = int64(0)
			num      = int64(0)
			afterNum = int64(0)
			source   = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.bagLogInfoc(uid, bagID, giftID, num, afterNum, source)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServicegiftActionInfoc(t *testing.T) {
	convey.Convey("giftActionInfoc", t, func(ctx convey.C) {
		var (
			uid      = int64(0)
			roomid   = int64(0)
			item     = int64(0)
			value    = int64(0)
			change   = int64(0)
			describe = ""
			platform = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.giftActionInfoc(uid, roomid, item, value, change, describe, platform)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServiceMakeID(t *testing.T) {
	convey.Convey("MakeID", t, func(ctx convey.C) {
		var (
			uid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := MakeID(uid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
