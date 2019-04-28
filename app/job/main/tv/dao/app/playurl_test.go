package app

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppPlayurl(t *testing.T) {
	var (
		c   = context.Background()
		cid = int(2578706)
	)
	convey.Convey("Playurl", t, func(ctx convey.C) {
		playurl, err := d.Playurl(c, cid)
		ctx.Convey("Then err should be nil.playurl should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(playurl, convey.ShouldNotBeNil)
		})
	})
}
