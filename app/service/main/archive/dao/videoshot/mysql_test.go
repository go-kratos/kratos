package videoshot

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestVideoshotvideoshot(t *testing.T) {
	var (
		c   = context.TODO()
		cid = int64(10097272)
	)
	convey.Convey("videoshot", t, func(ctx convey.C) {
		_, err := d.videoshot(c, cid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
