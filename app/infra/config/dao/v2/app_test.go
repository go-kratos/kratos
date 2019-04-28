package v2

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV2AppByTree(t *testing.T) {
	var (
		zone   = ""
		env    = ""
		treeID = int64(0)
	)
	convey.Convey("AppByTree", t, func(ctx convey.C) {
		app, err := d.AppByTree(zone, env, treeID)
		ctx.Convey("Then err should be nil.app should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(app, convey.ShouldNotBeNil)
		})
	})
}

func TestV2AppsByNameEnv(t *testing.T) {
	var (
		name = "main.common-arch.apm-admin"
		env  = "fat1"
	)
	convey.Convey("AppsByNameEnv", t, func(ctx convey.C) {
		apps, err := d.AppsByNameEnv(name, env)
		ctx.Convey("Then err should be nil.apps should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(apps, convey.ShouldNotBeNil)
		})
	})
}

func TestV2AppGet(t *testing.T) {
	var (
		zone  = "sh001"
		env   = "fat1"
		token = "a882c5530bcc11e8ab68522233017188"
	)
	convey.Convey("AppGet", t, func(ctx convey.C) {
		app, err := d.AppGet(zone, env, token)
		ctx.Convey("Then err should be nil.app should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(app, convey.ShouldNotBeNil)
		})
	})
}
