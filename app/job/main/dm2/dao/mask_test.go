package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaomaskURI(t *testing.T) {
	convey.Convey("maskURI", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.maskURI()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGenerateMask(t *testing.T) {
	convey.Convey("GenerateMask", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			cid      = int64(0)
			mid      = int64(0)
			plat     = int8(0)
			priority = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.GenerateMask(c, cid, mid, plat, priority, 0, 0, 0)
		})
	})
}
