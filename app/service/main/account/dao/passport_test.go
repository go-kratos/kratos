package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPassportDetail(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110016931)
	)
	convey.Convey("Get passport detail", t, func(ctx convey.C) {
		res, err := d.PassportDetail(c, mid)
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
