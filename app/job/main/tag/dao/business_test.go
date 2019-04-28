package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBusiness(t *testing.T) {
	var (
		c     = context.Background()
		state = int32(0)
	)
	convey.Convey("Business", t, func(ctx convey.C) {
		business, err := d.Business(c, state)
		ctx.Convey("Then err should be nil.business should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(business, convey.ShouldNotBeEmpty)
		})
	})
}
