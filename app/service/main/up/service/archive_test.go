package service

import (
	"context"
	"testing"

	upgrpc "go-common/app/service/main/up/api/v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpArcs(t *testing.T) {
	convey.Convey("UpArcs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upgrpc.UpArcsReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.UpArcs(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpsArcs(t *testing.T) {
	convey.Convey("UpsArcs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upgrpc.UpsArcsReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.UpsArcs(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpCount(t *testing.T) {
	convey.Convey("UpCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upgrpc.UpCountReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.UpCount(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpsCount(t *testing.T) {
	convey.Convey("UpsCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upgrpc.UpsCountReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.UpsCount(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpsAidPubTime(t *testing.T) {
	convey.Convey("UpsAidPubTime", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upgrpc.UpsArcsReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.UpsAidPubTime(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceAddUpPassedCache(t *testing.T) {
	convey.Convey("AddUpPassedCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upgrpc.UpCacheReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.AddUpPassedCache(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceDelUpPassedCache(t *testing.T) {
	convey.Convey("DelUpPassedCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &upgrpc.UpCacheReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.DelUpPassedCache(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
