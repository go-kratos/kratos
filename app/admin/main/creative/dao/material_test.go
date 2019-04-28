package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCategoryByID(t *testing.T) {
	convey.Convey("CategoryByID", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cate, err := d.CategoryByID(c, id)
			ctx.Convey("Then err should be nil.cate should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cate, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBindWithCategory(t *testing.T) {
	convey.Convey("BindWithCategory", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			MaterialID = int64(2)
			CategoryID = int64(1)
			index      = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.BindWithCategory(c, MaterialID, CategoryID, index)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}
