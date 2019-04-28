package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetByDiscountIds(t *testing.T) {
	convey.Convey("GetByDiscountIds", t, func(c convey.C) {
		var (
			ctx = context.Background()
			ids = []int64{1}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetByDiscountIds(ctx, ids)
			c.Convey("Then err should be nil.res should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
