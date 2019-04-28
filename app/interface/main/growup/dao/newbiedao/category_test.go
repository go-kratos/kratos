package newbiedao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewbiedaoGetCategories(t *testing.T) {
	convey.Convey("GetCategories", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.GetCategories(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
