package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoKeyAreas(t *testing.T) {
	convey.Convey("KeyAreas", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = "aid:2334"
			areas = []string{"reply"}
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			rs, err := d.KeyAreas(c, key, areas)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCoverStr(t *testing.T) {
	convey.Convey("CoverStr", t, func(ctx convey.C) {
		var (
			strs = []string{"str1", "str2"}
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			str := d.CoverStr(strs)
			ctx.Convey("Then str should equal 'strs[0]','strs[1]',...", func(ctx convey.C) {
				ctx.So(str, convey.ShouldEqual, "'str1','str2'")
			})
		})
	})
}
