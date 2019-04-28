package dao

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDaoparseCursor(t *testing.T) {
	convey.Convey("parseCursor", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
		)

		convCtx.Convey("common case", func(convCtx convey.C) {
			cursorPrev := ""
			cursorNext := ""
			cursor, directionNext, err := parseCursor(ctx, cursorPrev, cursorNext)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(directionNext, convey.ShouldBeTrue)
			convCtx.So(cursor.Offset, convey.ShouldEqual, 0)
			convCtx.So(cursor.StickRank, convey.ShouldEqual, 0)
		})

		convCtx.Convey("error case", func(convCtx convey.C) {
			cursorPrev := ""
			cursorNext := "{\"stick_rank\":1,\"offset\":1}"
			_, _, err := parseCursor(ctx, cursorPrev, cursorNext)
			convCtx.So(err, convey.ShouldNotBeNil)

			cursorPrev = "{\"stick_rank\":0,\"offset\":0}"
			cursorNext = ""
			_, _, err = parseCursor(ctx, cursorPrev, cursorNext)
			convCtx.So(err, convey.ShouldNotBeNil)

			cursorPrev = "{stick_rank\":0,\"offset\":0}"
			cursorNext = ""
			_, _, err = parseCursor(ctx, cursorPrev, cursorNext)
			convCtx.So(err, convey.ShouldNotBeNil)
		})

	})
}

func TestDaogetRedisList(t *testing.T) {
	convey.Convey("getRedisList", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			key = "stick:ttttt"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			list, err := d.getRedisList(ctx, key)
			convCtx.Convey("Then err should be nil.list should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(list, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaosetRedisList(t *testing.T) {
	convey.Convey("setRedisList", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			key  = "stick:topic"
			list = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.setRedisList(ctx, key, list)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
