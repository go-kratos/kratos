package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/videoup-task/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoListConfs(t *testing.T) {
	convey.Convey("ListConfs", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, err := d.ListConfs(c, []int64{}, "", "", "", 0, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReviewConfs(t *testing.T) {
	convey.Convey("ReviewConfs", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ReviewConfs(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpReviewConf(t *testing.T) {
	convey.Convey("UpReviewConf", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpReviewConf(c, &model.ReviewConf{})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelReviewConf(t *testing.T) {
	convey.Convey("DelReviewConf", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.DelReviewConf(c, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReviewForm(t *testing.T) {
	convey.Convey("ReviewForm", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ReviewForm(c, 0)
			ctx.Convey("Then err should be nil.max should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
