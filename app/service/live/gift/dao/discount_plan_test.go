package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetDiscountPlan(t *testing.T) {
	convey.Convey("GetDiscountPlan", t, func(c convey.C) {
		var (
			ctx = context.Background()
			now = time.Now()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			plans, err := d.GetDiscountPlan(ctx, now)
			c.Convey("Then err should be nil.plans should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(plans, convey.ShouldNotBeNil)
			})
		})
	})
}
