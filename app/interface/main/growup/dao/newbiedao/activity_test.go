package newbiedao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewbiedaoGetActivities(t *testing.T) {
	convey.Convey("GetActivities", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetActivities(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
