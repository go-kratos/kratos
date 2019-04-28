package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetOriginStreamingStatus(t *testing.T) {
	convey.Convey("SetOriginStreamingStatus", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rid  = int64(11891462)
			src  = int(1)
			from = int(1)
			to   = int(16)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetOriginStreamingStatus(c, rid, src, from, to)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
