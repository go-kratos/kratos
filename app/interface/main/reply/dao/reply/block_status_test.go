package reply

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyBlockInfo(t *testing.T) {
	convey.Convey("BlockInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.BlockStatus.BlockInfo(c, mid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
