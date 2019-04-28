package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetOnlinePlan(t *testing.T) {
	convey.Convey("GetOnlinePlan", t, func(c convey.C) {
		var (
			ctx = context.Background()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			plans, err := d.GetOnlinePlan(ctx)
			c.Convey("Then err should be nil.plans should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(plans, convey.ShouldNotBeNil)
			})
		})
	})
}
