package dao

import (
	"go-common/app/service/live/xuser/conf"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoNew(t *testing.T) {
	convey.Convey("New", t, func(ctx convey.C) {
		var (
			c = &conf.Config{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			dao := New(c)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInitAPI(t *testing.T) {
	convey.Convey("InitAPI", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			InitAPI()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDaogetConf(t *testing.T) {
	convey.Convey("getConf", t, func(ctx convey.C) {
		var (
			appName = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := getConf(appName)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
