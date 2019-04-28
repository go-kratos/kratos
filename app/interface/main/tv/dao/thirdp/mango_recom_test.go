package thirdp

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestThirdpMangoOrder(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("MangoOrder", t, func(ctx convey.C) {
		s, err := d.MangoOrder(c)
		ctx.Convey("Then err should be nil.s should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(s, convey.ShouldNotBeNil)
		})
	})
}

func TestThirdpMangoRecom(t *testing.T) {
	var (
		c   = context.Background()
		ids = []int64{}
	)
	convey.Convey("MangoRecom", t, func(ctx convey.C) {
		data, err := d.MangoRecom(c, ids)
		ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(data, convey.ShouldNotBeNil)
		})
	})
}
