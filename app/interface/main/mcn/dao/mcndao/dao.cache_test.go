package mcndao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaoMcnSign(t *testing.T) {
	convey.Convey("McnSign", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.McnSign(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoMcnDataSummary(t *testing.T) {
	convey.Convey("McnDataSummary", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			id           = int64(0)
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.McnDataSummary(c, id, generateDate)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}
