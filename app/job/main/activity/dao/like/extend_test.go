package like

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeAddExtend(t *testing.T) {
	convey.Convey("AddExtend", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "(13511,100)"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.AddExtend(c, query)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
