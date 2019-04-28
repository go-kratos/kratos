package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetAllGift(t *testing.T) {
	convey.Convey("GetAllGift", t, func(c convey.C) {
		var (
			ctx = context.Background()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			gifts, err := d.GetAllGift(ctx)
			c.Convey("Then err should be nil.gifts should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(gifts, convey.ShouldNotBeNil)
			})
		})
	})
}
