package up

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpGetHighAllyUps(t *testing.T) {
	convey.Convey("GetHighAllyUps", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetHighAllyUps(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeEmpty)
			})
		})
	})
}
