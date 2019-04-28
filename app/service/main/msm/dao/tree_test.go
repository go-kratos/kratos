package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotreeToken(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("treeToken", t, func(ctx convey.C) {
		token, err := d.treeToken(c)
		ctx.Convey("Then err should be nil.token should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(token, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTreeAppInfo(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("TreeAppInfo", t, func(ctx convey.C) {
		appInfo, err := d.TreeAppInfo(c)
		ctx.Convey("Then err should be nil.appInfo should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(appInfo, convey.ShouldNotBeNil)
		})
	})
}
