package service

import (
	"context"
	"go-common/app/service/bbq/user/api"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceAddLike(t *testing.T) {
	convey.Convey("AddLike", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			in = &api.LikeReq{Mid: 88895104, UpMid: 88895134, Opid: 2233}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.AddLike(c, in)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceCancelLike(t *testing.T) {
	convey.Convey("CancelLike", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			in = &api.LikeReq{Mid: 88895104, UpMid: 88895134, Opid: 2233}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.AddLike(c, &api.LikeReq{Mid: 88895104, UpMid: 88895134, Opid: 2233})
			res, err := s.CancelLike(c, in)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceListUserLike(t *testing.T) {
	convey.Convey("ListUserLike", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &api.ListUserLikeReq{UpMid: 88895104}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.ListUserLike(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceIsLike(t *testing.T) {
	convey.Convey("IsLike", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &api.IsLikeReq{Mid: 88895104, Svids: []int64{2233}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.AddLike(c, &api.LikeReq{Mid: 88895104, UpMid: 88895134, Opid: 2233})
			res, err := s.IsLike(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
