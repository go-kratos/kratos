package newbiedao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewbiedaoGetRecommendUpList(t *testing.T) {
	convey.Convey("GetRecommendUpList", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.GetRecommendUpList(c)
			ctx.Convey("Then err should be nil.recUps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
