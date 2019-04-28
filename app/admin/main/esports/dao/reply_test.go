package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRegReply(t *testing.T) {
	var (
		c    = context.Background()
		maid = int64(11)
		adid = int64(123)
	)
	convey.Convey("RegReply", t, func(ctx convey.C) {
		err := d.RegReply(c, maid, adid, "25")
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
