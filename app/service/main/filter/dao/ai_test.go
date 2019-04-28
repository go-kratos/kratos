package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAiWhite(t *testing.T) {
	convey.Convey("AiWhite", t, func(ctx convey.C) {
		var (
			c = context.TODO()
		)
		ctx.Convey("When everything looks good.", func(ctx convey.C) {
			res, err := d.AiWhite(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
