package abtest

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAbtestAbTest(t *testing.T) {
	convey.Convey("Abtest", t, func() {
		var (
			c      = context.Background()
			names  = "test"
			ipaddr = "0.0.0.0"
		)
		convey.Convey("When everything is correct", func(ctx convey.C) {
			adr, err := d.AbTest(context.TODO(), "", "")
			ctx.Convey("Then error should be nil，adr should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(adr, convey.ShouldNotBeNil)
			})
		})
		convey.Convey("When http request gets 404 error", func(ctx convey.C) {
			httpMock("GET", d.testURL).Reply(404)
			_, err := d.AbTest(c, names, ipaddr)
			ctx.Convey("Then error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		convey.Convey("When http request gets code != 0", func(ctx convey.C) {
			httpMock("GET", d.testURL).Reply(200).JSON(`{"code":-3,"message":"faild","data":{}}`)
			_, err := d.AbTest(c, names, ipaddr)
			ctx.Convey("Then error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
