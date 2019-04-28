package dao

import (
	"context"
	"go-common/app/job/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaospecialDmKey(t *testing.T) {
	convey.Convey("specialDmKey", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := testDao.specialDmKey(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelSpecialDmCache(t *testing.T) {
	convey.Convey("DelSpecialDmCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			testDao.DelSpecialDmCache(c, oid, tp)
		})
	})
}

func TestDaoAddSpecialDmCache(t *testing.T) {
	convey.Convey("AddSpecialDmCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ds = &model.DmSpecial{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := testDao.AddSpecialDmCache(c, ds)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
