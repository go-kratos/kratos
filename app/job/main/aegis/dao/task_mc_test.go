package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoIsConsumerOn(t *testing.T) {
	convey.Convey("IsConsumerOn", t, func(ctx convey.C) {
	})
}

func TestDaomcKey(t *testing.T) {
	convey.Convey("mcKey", t, func(ctx convey.C) {
		var (
			bizid  = int(0)
			flowid = int(0)
			uid    = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mcKey(bizid, flowid, uid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
