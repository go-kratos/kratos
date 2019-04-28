package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAllAppsInfo(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("AllAppsInfo", t, func(ctx convey.C) {
		res, err := d.AllAppsInfo(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAllAppsAuth(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("AllAppsAuth", t, func(ctx convey.C) {
		res, err := d.AllAppsAuth(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
