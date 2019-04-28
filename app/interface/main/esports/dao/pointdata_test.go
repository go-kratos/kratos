package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLeida(t *testing.T) {
	convey.Convey("Leida", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			url = "http://47.95.28.113/nesport/index.php/Api/lol/games?match_id=328266&key=d076eef519773c5954081e6a352c726d"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rs, err := d.Leida(c, url)
			convCtx.Convey("Then err should be nil.rs should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoThirdGet(t *testing.T) {
	convey.Convey("ThirdGet", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			url = "http://47.95.28.113/nesport/index.php/Api/lol/games?match_id=328266&key=d076eef519773c5954081e6a352c726d"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.ThirdGet(c, url)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
