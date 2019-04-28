package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveArcsAids(t *testing.T) {
	convey.Convey("ArcsAids", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{10110188}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			aids, ptimes, copyrights, aptm, err := d.ArcsAids(c, ids)
			ctx.Convey("Then err should be nil.aids,ptimes,copyrights,aptm should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(aptm, convey.ShouldNotBeNil)
				ctx.So(copyrights, convey.ShouldNotBeNil)
				ctx.So(ptimes, convey.ShouldNotBeNil)
				ctx.So(aids, convey.ShouldNotBeNil)
			})
		})
	})
}
