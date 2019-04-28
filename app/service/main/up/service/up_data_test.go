package service

import (
	"context"
	"testing"

	upgrpc "go-common/app/service/main/up/api/v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpBaseStats(t *testing.T) {
	convey.Convey("UpBaseStats", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			req = &upgrpc.UpStatReq{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpBaseStats(c, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
