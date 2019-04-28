package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAutoCaseConf(t *testing.T) {
	convey.Convey("AutoCaseConf", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			otype = int8(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ac, err := d.AutoCaseConf(c, otype)
			ctx.Convey("Then err should be nil.ac should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ac, convey.ShouldNotBeNil)
			})
		})
	})
}
