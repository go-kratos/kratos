package v1

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV1bagLogInfoc(t *testing.T) {
	convey.Convey("bagLogInfoc", t, func(ctx convey.C) {
		var (
			uid      = int64(5)
			bagID    = int64(10)
			giftID   = int64(1)
			num      = int64(1)
			afterNum = int64(2)
			source   = "xx"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.bagLogInfoc(uid, bagID, giftID, num, afterNum, source)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestV1MakeID(t *testing.T) {
	convey.Convey("MakeID", t, func(ctx convey.C) {
		var (
			uid = int64(5)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := MakeID(uid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
				ctx.So(p1, convey.ShouldContainSubstring, "50000000005")
			})
		})
	})
}

func TestV1infoc(t *testing.T) {
	convey.Convey("infoc", t, func(ctx convey.C) {
		var (
			i = interface{}(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.infoc(i)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestV1infocproc(t *testing.T) {
	convey.Convey("infocproc", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.infocproc()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
