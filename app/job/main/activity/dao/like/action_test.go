package like

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeBatchLikeActSum(t *testing.T) {
	convey.Convey("BatchLikeActSum", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			lids = []int64{13511, 13512, 13510}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.BatchLikeActSum(c, lids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
