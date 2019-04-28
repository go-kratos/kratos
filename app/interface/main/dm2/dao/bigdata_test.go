package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaodanmakuURI(t *testing.T) {
	convey.Convey("danmakuURI", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.danmakuURI(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRecFlag(t *testing.T) {
	convey.Convey("RecFlag", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			aid   = int64(0)
			oid   = int64(0)
			limit = int64(0)
			ps    = int64(0)
			pe    = int64(0)
			plat  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.RecFlag(c, mid, aid, oid, limit, ps, pe, plat)
		})
	})
}
