package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceArcTopDataStatistics(t *testing.T) {
	convey.Convey("ArcTopDataStatistics", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnGetRankReq{}
		)
		arg.SignID = 214
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.ArcTopDataStatistics(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnsTotalDatas(t *testing.T) {
	convey.Convey("McnsTotalDatas", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.TotalMcnDataReq{Date: 1542211200}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.McnsTotalDatas(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicestatTypeRate(t *testing.T) {
	convey.Convey("statTypeRate", t, func(ctx convey.C) {
		var (
			dts   = []*model.DataTypes{}
			total = int64(111)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			statTypeRate(dts, total)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
