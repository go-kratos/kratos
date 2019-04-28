package subtitle

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSubtitleUpdate(t *testing.T) {
	convey.Convey("Update", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aid  = int64(10110826)
			open bool
			lan  = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Update(c, aid, open, lan)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
