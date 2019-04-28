package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/apm/model/ut"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceDashCurveGraph(t *testing.T) {
	convey.Convey("ProjectCurveGraph", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &ut.PCurveReq{
				StartTime: 1536508800,
				EndTime:   1541779200,
			}
			username = "fengshanshan"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := svr.DashCurveGraph(c, username, req)
			for _, r := range res {
				t.Logf("res:%+v", r)
			}
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceDashGraphDetail(t *testing.T) {
	convey.Convey("ProjectGraphDetail", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			username = "haoguanwei"
			req      = &ut.PCurveReq{
				StartTime: 1536508800,
				EndTime:   1541779200,
				Path:      "go-common/app/admin/main/apm",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := svr.DashGraphDetail(c, username, req)
			for _, r := range res {
				t.Logf("res:%+v", r)
			}
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceDashGraphDetailSingle(t *testing.T) {
	convey.Convey("ProjectGraphDetailSingle", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			username = "haoguanwei"
			req      = &ut.PCurveReq{
				User:      "fengshanshan",
				StartTime: 1536508800,
				EndTime:   1541779200,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := svr.DashGraphDetailSingle(c, username, req)
			for _, r := range res {
				t.Logf("res:%+v\n", r)
			}
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceDashPkgs(t *testing.T) {
	convey.Convey("DashboardPkgs", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			path     = "go-common/app/"
			username = "zhaobingqing"
		)
		ctx.Convey("When path is none", func(ctx convey.C) {
			val, err := svr.DashPkgsTree(c, path, username)
			ctx.Convey("Error should be nil, pkgs should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(val, convey.ShouldNotBeNil)
				t.Logf("the val are %+v", val)
			})
		})
		// ctx.Convey("When path is not none", func(ctx convey.C) {
		// 	path = "go-common/app/service/main/block"
		// 	val, err := svr.GetPersonalPkgs(c, path, username)
		// 	ctx.Convey("Error should be nil, pkgs should not be nil", func(ctx convey.C) {
		// 		ctx.So(err, convey.ShouldBeNil)
		// 		ctx.So(val, convey.ShouldNotBeNil)
		// 		t.Logf("the val are %+v", val)
		// 	})
		// })
	})
}
func TestServiceAppsCache(t *testing.T) {
	convey.Convey("AppsCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := svr.AppsCache(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
