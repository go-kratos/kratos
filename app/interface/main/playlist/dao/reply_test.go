package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRegReply(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
		mid = int64(2)
	)
	convey.Convey("RegReply", t, func(ctx convey.C) {
		err := d.RegReply(c, pid, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
