package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceRecommendAdd(t *testing.T) {
	convey.Convey("RecommendAdd", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.RecommendUpReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.RecommendAdd(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceRecommendOP(t *testing.T) {
	convey.Convey("RecommendOP", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.RecommendStateOpReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.RecommendOP(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceRecommendList(t *testing.T) {
	convey.Convey("RecommendList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPRecommendReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.RecommendList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
