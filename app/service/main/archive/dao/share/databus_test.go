package share

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSharePubStatDatabus(t *testing.T) {
	var (
		c     = context.TODO()
		aid   = int64(1)
		share = int(1)
	)
	convey.Convey("PubStatDatabus", t, func(ctx convey.C) {
		err := d.PubStatDatabus(c, aid, share)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestSharePubShare(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		mid = int64(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("PubShare", t, func(ctx convey.C) {
		err := d.PubShare(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
