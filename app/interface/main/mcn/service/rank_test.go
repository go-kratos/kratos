package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/mcn/model/mcnmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceMcnGetRankUpFans(t *testing.T) {
	convey.Convey("McnGetRankUpFans", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnGetRankReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnGetRankUpFans(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnGetRankArchiveLikes(t *testing.T) {
	convey.Convey("McnGetRankArchiveLikes", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnGetRankReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnGetRankArchiveLikes(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
