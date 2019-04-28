package like

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeRawActLikeAchieves(t *testing.T) {
	convey.Convey("RawActLikeAchieves", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10805)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.RawActLikeAchieves(c, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
