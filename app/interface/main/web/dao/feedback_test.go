package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/web/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoFeedback(t *testing.T) {
	convey.Convey("Feedback", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			feedParams = &model.Feedback{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Feedback(c, feedParams)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
