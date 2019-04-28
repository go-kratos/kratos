package app

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppPgcSeaSug(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("PgcSeaSug", t, func(ctx convey.C) {
		res, err := d.PgcSeaSug(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestAppUgcSeaSug(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("UgcSeaSug", t, func(ctx convey.C) {
		res, err := d.UgcSeaSug(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
