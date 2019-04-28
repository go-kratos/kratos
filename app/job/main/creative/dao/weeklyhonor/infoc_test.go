package weeklyhonor

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWeeklyhonorHonorInfoc(t *testing.T) {
	convey.Convey("HonorInfoc", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			success = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.HonorInfoc(c, mid, success)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
