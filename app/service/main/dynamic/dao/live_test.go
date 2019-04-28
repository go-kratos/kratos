package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLive(t *testing.T) {
	convey.Convey("Live", t, func() {
		var (
			c = context.Background()
		)
		convey.Convey("When http request gets code == 0", func(ctx convey.C) {
			httpMock("GET", d.liveURI).Reply(200).JSON(`{"code":0,"msg":"ok","message":"ok","data":{"count":770843}}`)
			count, err := d.Live(c)
			ctx.Convey("Then err should be nil.count should greater 0.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}
