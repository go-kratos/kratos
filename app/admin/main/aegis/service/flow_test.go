package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/net"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceShowFlow(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("ShowFlow", t, func(ctx convey.C) {
		r, err := s.ShowFlow(c, id)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceGetFlowList(t *testing.T) {
	var (
		c  = context.TODO()
		pm = &net.ListNetElementParam{}
	)
	convey.Convey("GetFlowList", t, func(ctx convey.C) {
		result, err := s.GetFlowList(c, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestServicecheckFlowUnique(t *testing.T) {
	var (
		netID = int64(0)
		name  = ""
	)
	convey.Convey("checkFlowUnique", t, func(ctx convey.C) {
		err, msg := s.checkFlowUnique(cntx, netID, name)
		ctx.Convey("Then err should be nil.msg should not be nil.", func(ctx convey.C) {
			ctx.So(msg, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicetaskFlowByBusiness(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("dispatchFlow", t, func(ctx convey.C) {
		res, err := s.dispatchFlow(c, []int64{1, 2}, []int64{1, 2})
		for biz, item := range res {
			t.Logf("biz(%d) flows(%+v)", biz, item)
		}
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
