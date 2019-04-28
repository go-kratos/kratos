package report

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReportSetReportCache(t *testing.T) {
	convey.Convey("SetReportCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = make(map[string]interface{})
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetReportCache(c, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestReportGetReportCache(t *testing.T) {
	convey.Convey("GetReportCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetReportCache(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
