package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaocacheSFEquip(t *testing.T) {
	convey.Convey("cacheSFEquip", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.cacheSFEquip(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
