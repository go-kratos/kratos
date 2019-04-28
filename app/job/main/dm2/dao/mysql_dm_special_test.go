package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpsertDmSpecialLocation(t *testing.T) {
	convey.Convey("UpsertDmSpecialLocation", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tp        = int32(0)
			oid       = int64(0)
			locations = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := testDao.UpsertDmSpecialLocation(c, tp, oid, locations)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDMSpecialLocations(t *testing.T) {
	convey.Convey("DMSpecialLocations", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(1)
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DMSpecialLocations(c, tp, oid)
		})
	})
}
