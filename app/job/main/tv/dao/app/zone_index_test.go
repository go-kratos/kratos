package app

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppPassedSn(t *testing.T) {
	var (
		c        = context.Background()
		category = int(1)
	)
	convey.Convey("PassedSn", t, func(ctx convey.C) {
		res, err := d.PassedSn(c, category)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
