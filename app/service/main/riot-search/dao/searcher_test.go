package dao

import (
	"go-common/app/service/main/riot-search/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSearchIDOnly(t *testing.T) {
	convey.Convey("SearchIDOnly", t, func(ctx convey.C) {
		var (
			arg1 = &model.RiotSearchReq{}
			arg2 = &model.RiotSearchReq{Keyword: "test", IDs: []uint64{1}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.SearchIDOnly(arg1)
			ctx.Convey("Then p1 should be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldBeNil)
			})
			p2 := d.SearchIDOnly(arg2)
			ctx.Convey("Then p2 should not be nil.", func(ctx convey.C) {
				ctx.So(p2, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSearch(t *testing.T) {
	convey.Convey("Search", t, func(ctx convey.C) {
		var (
			arg1 = &model.RiotSearchReq{}
			arg2 = &model.RiotSearchReq{Keyword: "test", IDs: []uint64{1}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.Search(arg1)
			ctx.Convey("Then p1 should be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldBeNil)
			})
			p2 := d.Search(arg2)
			ctx.Convey("Then p2 should not be nil.", func(ctx convey.C) {
				ctx.So(p2, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHas(t *testing.T) {
	convey.Convey("Search", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.Has(1)
			ctx.Convey("Then p1 should be false.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldBeFalse)
			})
		})
	})
}
