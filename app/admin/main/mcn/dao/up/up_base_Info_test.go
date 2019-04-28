package up

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpUpBaseInfoMap(t *testing.T) {
	convey.Convey("UpBaseInfoMap", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mbi, err := d.UpBaseInfoMap(c, mids)
			ctx.Convey("Then err should be nil.mbi should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mbi, convey.ShouldNotBeNil)
			})
		})
	})
}
