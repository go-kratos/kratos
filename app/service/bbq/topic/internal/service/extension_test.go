package service

import (
	"context"
	"go-common/app/service/bbq/topic/api"
	"go-common/library/log"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceRegister(t *testing.T) {
	convey.Convey("Register", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			//req = &api.VideoExtension{Svid: 1, Extension: "{\"title_extra\":[{\"end\":5,\"topicId\":0,\"type\":1,\"name\":\"#Test\",\"start\":0}]}"}
			req = &api.VideoExtension{Svid: 1, Extension: "{\"title_extra\":[{\"name\":\"Test\",\"start\":0, \"end\":4},{\"name\":\"test_333333\",\"start\":10, \"end\":18}]}"}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.Register(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceListExtension(t *testing.T) {
	convey.Convey("ListExtension", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			req = &api.ListExtensionReq{Svids: []int64{1}}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.ListExtension(ctx, req)
			log.V(1).Infow(ctx, "res", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
