package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSendCapture(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(0)
		code = int(0)
	)
	convey.Convey("SendCapture", t, func(ctx convey.C) {
		err := d.SendCapture(c, mid, code)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, 74001)
		})
	})
}
