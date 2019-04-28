package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetAllParameter(t *testing.T) {
	convey.Convey("GetAllParameter", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			parameters, err := d.GetAllParameter(c)
			ctx.Convey("Then err should be nil.parameters should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(parameters, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertParameter(t *testing.T) {
	convey.Convey("InsertParameter", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = "('wmaafans',100,'活跃粉丝数')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertParameter(c, val)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
