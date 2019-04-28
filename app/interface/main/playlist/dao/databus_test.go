package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPubView(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		aid  = int64(13825646)
		view = int64(1000)
	)
	convey.Convey("PubView", t, func(ctx convey.C) {
		err := d.PubView(c, pid, aid, view)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPubShare(t *testing.T) {
	var (
		c     = context.Background()
		pid   = int64(1)
		aid   = int64(13825646)
		share = int64(50)
	)
	convey.Convey("PubShare", t, func(ctx convey.C) {
		err := d.PubShare(c, pid, aid, share)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
