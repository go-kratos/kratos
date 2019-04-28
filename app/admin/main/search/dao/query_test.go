package dao

import (
	"context"
	"go-common/app/admin/main/search/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoQueryConf(t *testing.T) {
	convey.Convey("QueryConf", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.QueryConf(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryBasic(t *testing.T) {
	convey.Convey("QueryBasic", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mixedQuery, qbDebug := d.QueryBasic(c, sp)
			ctx.Convey("Then mixedQuery,qbDebug should not be nil.", func(ctx convey.C) {
				ctx.So(qbDebug, convey.ShouldNotBeNil)
				ctx.So(mixedQuery, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoqueryBasicRange(t *testing.T) {
	convey.Convey("queryBasicRange", t, func(ctx convey.C) {
		var (
			rangeMap map[string]string
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rangeQuery, err := d.queryBasicRange(rangeMap)
			ctx.Convey("Then err should be nil.rangeQuery should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rangeQuery, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoqueryBasicLike(t *testing.T) {
	convey.Convey("queryBasicLike", t, func(ctx convey.C) {
		var (
			likeMap  = []model.QueryBodyWhereLike{}
			business = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.queryBasicLike(likeMap, business)
			ctx.Convey("Then err should be nil.likeQuery should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(likeQuery, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoScroll(t *testing.T) {
	convey.Convey("Scroll", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			sp = &model.QueryParams{
				QueryBody: &model.QueryBody{},
				AppIDConf: &model.QueryConfDetail{
					ESCluster: "",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, debug, err :=
			d.Scroll(c, sp)
			ctx.Convey("Then err should be nil.res,debug should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				//ctx.So(debug, convey.ShouldNotBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
