package web

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWebPushAll(t *testing.T) {
	var (
		c   = context.Background()
		msg = `{"1":64060766,"180":23753334,"31":49310729,"live":516587}`
		ip  = "127.0.0.1"
	)
	convey.Convey("PushAll", t, func(ctx convey.C) {
		err := d.PushAll(c, msg, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
