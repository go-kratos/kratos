package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPlatformSuit(t *testing.T) {
	convey.Convey("PlatformSuit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ConfPlatform{
				ID:           1,
				PlatformName: "android-ut",
				Platform:     "android",
				Device:       "phone",
				MobiApp:      "android",
				PanelType:    "normal",
				Operator:     "test1",
			}
		)
		ctx.Convey("clean data before", func(ctx convey.C) {
			d.vip.Table(_vipConfPlatform).Where("platform_name=?", "android-ut").Delete(model.ConfPlatform{})
		})
		ctx.Convey("PlatformSave", func(ctx convey.C) {
			eff, err := d.PlatformSave(c, arg)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(eff, convey.ShouldNotBeNil)
		})
		ctx.Convey("PlatformByID", func(ctx convey.C) {
			re, err := d.PlatformByID(c, arg.ID)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(re, convey.ShouldNotBeNil)
		})
		ctx.Convey("PlatformAll", func(ctx convey.C) {
			res, err := d.PlatformAll(c, "desc")
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		ctx.Convey("PlatformDel", func(ctx convey.C) {
			eff, err := d.PlatformDel(c, arg.ID, arg.Operator)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(eff, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		ctx.Convey("clean data", func(ctx convey.C) {
			d.vip.Table(_vipConfPlatform).Where("platform_name=?", "android-ut").Delete(model.ConfPlatform{})
		})
	})
}

func TestDaoPlatformTypes(t *testing.T) {
	convey.Convey("PlatformTypes", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.PlatformTypes(c)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
