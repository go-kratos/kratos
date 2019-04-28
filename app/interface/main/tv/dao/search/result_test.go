package search

import (
	"context"
	searchMdl "go-common/app/interface/main/tv/model/search"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchSearchAll(t *testing.T) {
	var (
		ctx = context.Background()
		req = &searchMdl.ReqSearch{
			SearchType: "all_tv",
			Keyword:    "番剧",
			Page:       1,
			MobiAPP:    "android_tv_yst",
			Platform:   "android",
			Build:      "1011",
		}
	)
	convey.Convey("SearchAll", t, func(cx convey.C) {
		result, common, err := d.SearchAll(ctx, req)
		cx.Convey("Then err should be nil.result,all should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(common, convey.ShouldNotBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestSearchSearchUgc(t *testing.T) {
	var (
		ctx = context.Background()
		req = &searchMdl.ReqSearch{
			SearchType: "tv_ugc",
			Category:   160,
			Keyword:    "测试",
			Page:       1,
			MobiAPP:    "android_tv_yst",
			Platform:   "android",
			Build:      "1011",
		}
	)
	convey.Convey("SearchUgc", t, func(cx convey.C) {
		result, common, err := d.SearchUgc(ctx, req)
		cx.Convey("Then err should be nil.result,ugc should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(common, convey.ShouldNotBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestSearchSearchPgc(t *testing.T) {
	var (
		ctx = context.Background()
		req = &searchMdl.ReqSearch{
			SearchType: "tv_pgc",
			Keyword:    "番剧",
			Page:       1,
			MobiAPP:    "android_tv_yst",
			Platform:   "android",
			Build:      "1011",
		}
	)
	convey.Convey("SearchPgc", t, func(cx convey.C) {
		result, common, err := d.SearchPgc(ctx, req)
		cx.Convey("Then err should be nil.result,pgc should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(common, convey.ShouldNotBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestSearchcommonParam(t *testing.T) {
	var (
		req = &searchMdl.ReqSearch{}
	)
	convey.Convey("commonParam", t, func(cx convey.C) {
		params := commonParam(req)
		cx.Convey("Then params should not be nil.", func(ctx convey.C) {
			ctx.So(params, convey.ShouldNotBeNil)
		})
	})
}
