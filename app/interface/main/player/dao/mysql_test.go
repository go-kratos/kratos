package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoParam(t *testing.T) {
	convey.Convey("Param", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		time.Sleep(time.Second)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			param, err := d.Param(c)
			ctx.Convey("Then err should be nil.param should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(param, convey.ShouldNotBeNil)
			})
		})
	})
}
