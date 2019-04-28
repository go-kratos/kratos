package service

import (
	"context"
	"go-common/app/service/main/up/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceGetUpInfoActive(t *testing.T) {
	convey.Convey("GetUpInfoActive", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.UpInfoActiveReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.GetUpInfoActive(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
