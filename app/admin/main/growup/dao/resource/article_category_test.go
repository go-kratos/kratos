package resource

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestResourcecolumnCategory(t *testing.T) {
	convey.Convey("columnCategory", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := columnCategory(c)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestResourceColumnCategoryNameToID(t *testing.T) {
	convey.Convey("ColumnCategoryNameToID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			categories, err := ColumnCategoryNameToID(c)
			ctx.Convey("Then err should be nil.categories should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(categories, convey.ShouldNotBeNil)
			})
		})
	})
}
