package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetToken(t *testing.T) {
	var (
		c   = context.Background()
		bid = "account"
	)
	convey.Convey("GetToken", t, func(ctx convey.C) {
		res, err := d.GetToken(c, bid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
