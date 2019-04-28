package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddMcnUpRecommend(t *testing.T) {
	convey.Convey("AddMcnUpRecommend", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnUpRecommendPool{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.AddMcnUpRecommend(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelMcnUpRecommendPool(t *testing.T) {
	convey.Convey("DelMcnUpRecommendPool", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelMcnUpRecommendPool(c)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelMcnUpRecommendSource(t *testing.T) {
	convey.Convey("DelMcnUpRecommendSource", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelMcnUpRecommendSource(c, id)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMcnUpRecommendSources(t *testing.T) {
	convey.Convey("McnUpRecommendSources", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			limit = 100
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rps, err := d.McnUpRecommendSources(c, limit)
			ctx.Convey("Then err should be nil.rps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(rps), convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}
