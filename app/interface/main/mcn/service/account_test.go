package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/mcn/model/mcnmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceGetUpAccountInfo(t *testing.T) {
	convey.Convey("GetUpAccountInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.McnGetAccountReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetUpAccountInfo(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
