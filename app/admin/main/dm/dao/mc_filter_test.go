package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyUpFilter(t *testing.T) {
	convey.Convey("keyUpFilter", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyUpFilter(mid, oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelUpFilterCache(t *testing.T) {
	convey.Convey("DelUpFilterCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.DelUpFilterCache(c, mid, oid)
		})
	})
}
