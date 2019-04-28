package manager

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerSpecials(t *testing.T) {
	convey.Convey("Specials", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sps, err := d.Specials(c)
			ctx.Convey("Then err should be nil.sps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sps, convey.ShouldNotBeNil)
			})
		})
	})
}
