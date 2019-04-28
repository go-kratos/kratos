package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPubLabour(t *testing.T) {
	convey.Convey("PubLabour", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
			msg = interface{}(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.PubLabour(c, aid, msg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
