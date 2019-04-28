package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAuth(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Auth", t, func(ctx convey.C) {
		auths, err := d.Auth(c)
		ctx.Convey("Then err should be nil.auths should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(auths, convey.ShouldNotBeNil)
		})
	})
}
