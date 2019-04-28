package service

import (
	"context"
	"testing"

	"go-common/app/admin/ep/saga/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_QueryProjectRunners(t *testing.T) {
	convey.Convey("get saga-admin runners", t, func() {
		var (
			ctx = context.Background()
			req = &model.ProjectDataReq{ProjectID: 682}
		)

		resp, err := srv.QueryProjectRunners(ctx, req)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(resp), convey.ShouldBeGreaterThan, 1)
	})
	convey.Convey("get saga-admin runners", t, func() {
		var (
			ctx = context.Background()
			req = &model.ProjectDataReq{ProjectID: 4928}
		)

		resp, err := srv.QueryProjectRunners(ctx, req)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(resp), convey.ShouldBeGreaterThan, 1)
	})
	convey.Convey("get saga-admin runners", t, func() {
		var (
			ctx = context.Background()
			req = &model.ProjectDataReq{ProjectID: 5822}
		)

		resp, err := srv.QueryProjectRunners(ctx, req)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(resp), convey.ShouldBeGreaterThan, 1)
	})
}
