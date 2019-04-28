package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeGetBubbleMeta(t *testing.T) {
	convey.Convey("GetBubbleMeta", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int64(20)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, last, err := d.GetBubbleMeta(c, id, limit)
			ctx.Convey("Then err should be nil.data,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}
