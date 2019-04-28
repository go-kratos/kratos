package manager

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerUpSpecial(t *testing.T) {
	convey.Convey("UpSpecial", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(27515256)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.UpSpecial(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestManagerUpsSpecial(t *testing.T) {
	convey.Convey("UpsSpecial", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{27515256}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.UpsSpecial(c, keys)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
