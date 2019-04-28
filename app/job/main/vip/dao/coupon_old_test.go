package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSalaryLogMaxID(t *testing.T) {
	convey.Convey("SalaryLogMaxID", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			dv = "2018_09"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			maxID, err := d.SalaryLogMaxID(c, dv)
			ctx.Convey("Then err should be nil.maxID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSelOldSalaryList(t *testing.T) {
	convey.Convey("SelOldSalaryList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int(0)
			endID = int(0)
			dv    = "2018_09"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SelOldSalaryList(c, id, endID, dv)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
