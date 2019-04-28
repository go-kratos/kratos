package watermark

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWatermarkGenWm(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(0)
		wmText = "text"
		ip     = ""
	)
	convey.Convey("GenWm", t, func(ctx convey.C) {
		gm, err := d.GenWm(c, mid, wmText, ip)
		ctx.Convey("Then err should be nil.gm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(gm, convey.ShouldNotBeNil)
		})
	})
}
